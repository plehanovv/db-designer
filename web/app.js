const textArea = document.querySelector("#domainText");
const button = document.querySelector("#analyzeButton");
const statusNode = document.querySelector("#status");
const entitiesNode = document.querySelector("#entities");
const sqlNode = document.querySelector("#sql");
const diagram = document.querySelector("#diagram");

button.addEventListener("click", analyze);
window.addEventListener("load", analyze);

async function analyze() {
    const text = textArea.value.trim();
    if (!text) {
        statusNode.textContent = "Domain description is empty.";
        return;
    }

    button.disabled = true;
    statusNode.textContent = "Running semantic analysis...";

    try {
        const response = await fetch("/analyze", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({text}),
        });

        if (!response.ok) {
            throw new Error(await response.text());
        }

        const result = await response.json();
        renderEntities(result.entities || []);
        renderSQL(result.sql || "");
        renderDiagram(result.entities || [], result.relations || []);
        statusNode.textContent = `Entities: ${(result.entities || []).length}, relations: ${(result.relations || []).length}.`;
    } catch (error) {
        statusNode.textContent = `Analysis error: ${error.message}`;
    } finally {
        button.disabled = false;
    }
}

function renderEntities(entities) {
    entitiesNode.innerHTML = "";

    for (const entity of entities) {
        const card = document.createElement("article");
        card.className = "entity";
        card.innerHTML = `<div class="entity-name">${escapeHTML(entity.name)}</div>`;

        const attributes = entity.attributes?.length ? entity.attributes : [{name: "id", type: "SERIAL"}];
        for (const attr of attributes) {
            const row = document.createElement("div");
            row.className = "attribute";
            row.innerHTML = `<span>${escapeHTML(attr.name)}</span><span>${escapeHTML(attr.type)}</span>`;
            card.append(row);
        }

        entitiesNode.append(card);
    }
}

function renderSQL(sql) {
    sqlNode.textContent = sql || "-- SQL will appear after analysis";
}

function renderDiagram(entities, relations) {
    const width = Math.max(diagram.clientWidth, 720);
    const height = Math.max(diagram.clientHeight, 320);
    diagram.setAttribute("viewBox", `0 0 ${width} ${height}`);
    diagram.innerHTML = "";

    const positions = new Map();
    const radius = Math.min(width, height) * 0.34;
    const centerX = width / 2;
    const centerY = height / 2;

    entities.forEach((entity, index) => {
        const angle = (Math.PI * 2 * index) / Math.max(entities.length, 1) - Math.PI / 2;
        positions.set(entity.name, {
            x: centerX + Math.cos(angle) * radius,
            y: centerY + Math.sin(angle) * radius,
        });
    });

    for (const relation of relations) {
        const from = positions.get(relation.from);
        const to = positions.get(relation.to);
        if (!from || !to) continue;

        line(from.x, from.y, to.x, to.y, "#94a3b8");
        text((from.x + to.x) / 2, (from.y + to.y) / 2 - 8, relation.cardinality || relation.type, "12px", "#475569");
    }

    for (const entity of entities) {
        const position = positions.get(entity.name);
        rect(position.x - 82, position.y - 28, 164, 56, "#ffffff", "#2563eb");
        text(position.x, position.y - 4, entity.name, "14px", "#1e293b", "700");
        text(position.x, position.y + 16, `${entity.attributes?.length || 0} attrs`, "12px", "#64748b");
    }
}

function rect(x, y, width, height, fill, stroke) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    node.setAttribute("x", x);
    node.setAttribute("y", y);
    node.setAttribute("width", width);
    node.setAttribute("height", height);
    node.setAttribute("rx", "8");
    node.setAttribute("fill", fill);
    node.setAttribute("stroke", stroke);
    node.setAttribute("stroke-width", "1.5");
    diagram.append(node);
}

function line(x1, y1, x2, y2, stroke) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "line");
    node.setAttribute("x1", x1);
    node.setAttribute("y1", y1);
    node.setAttribute("x2", x2);
    node.setAttribute("y2", y2);
    node.setAttribute("stroke", stroke);
    node.setAttribute("stroke-width", "1.5");
    diagram.append(node);
}

function text(x, y, value, size, fill, weight = "500") {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "text");
    node.setAttribute("x", x);
    node.setAttribute("y", y);
    node.setAttribute("text-anchor", "middle");
    node.setAttribute("font-size", size);
    node.setAttribute("font-weight", weight);
    node.setAttribute("fill", fill);
    node.textContent = value;
    diagram.append(node);
}

function escapeHTML(value) {
    return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;");
}
