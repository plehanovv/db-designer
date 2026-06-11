const textArea = document.querySelector("#domainText");
const button = document.querySelector("#analyzeButton");
const uploadFileButton = document.querySelector("#uploadFileButton");
const sourceFileInput = document.querySelector("#sourceFile");
const addEntityButton = document.querySelector("#addEntityButton");
const addRelationButton = document.querySelector("#addRelationButton");
const regenerateButton = document.querySelector("#regenerateButton");
const downloadSQLButton = document.querySelector("#downloadSQLButton");
const downloadJSONButton = document.querySelector("#downloadJSONButton");
const downloadMermaidButton = document.querySelector("#downloadMermaidButton");
const downloadReportButton = document.querySelector("#downloadReportButton");
const databaseNameInput = document.querySelector("#databaseName");
const statusNode = document.querySelector("#status");
const entitiesNode = document.querySelector("#entities");
const relationsNode = document.querySelector("#relations");
const diagnosticsNode = document.querySelector("#diagnostics");
const traceNode = document.querySelector("#trace");
const transformationsNode = document.querySelector("#transformations");
const evaluationNode = document.querySelector("#evaluation");
const sqlNode = document.querySelector("#sql");
const diagramViewport = document.querySelector("#diagramViewport");
const diagram = document.querySelector("#diagram");
const diagramZoomOutButton = document.querySelector("#diagramZoomOutButton");
const diagramZoomResetButton = document.querySelector("#diagramZoomResetButton");
const diagramZoomInButton = document.querySelector("#diagramZoomInButton");
const diagramFitButton = document.querySelector("#diagramFitButton");
let currentModel = {database: {name: "database"}, entities: [], relations: []};
let currentAnalysis = {input: "", diagnostics: [], transformations: [], explanation: {candidates: []}, sql: ""};
let modelDirty = false;
let diagramPositions = new Map();
let activeDrag = null;
let activePan = null;
let diagramZoom = 1;
let diagramBaseSize = {width: 1800, height: 1100};
let diagramNeedsCenter = false;

const diagramNodeMargin = {x: 190, y: 130};
const diagramComponentMargin = {x: 460, y: 330};

let examples = {
    "library": "Читатель имеет имя email телефон.\nАвтор имеет имя страну.\nКнига имеет название год isbn.\nКатегория имеет название код.\nВыдача имеет дату выдачи дату возврата статус.\nБронирование имеет дату статус.\nШтраф имеет сумму дату статус.\nАвтор пишет книги.\nКатегория включает книги.\nЧитатель имеет несколько выдач.\nВыдача связана с книгой.\nЧитатель бронирует книги.\nБронирование связано с книгой.\nВыдача содержит штрафы.",
    "university": "Студент имеет имя email.\nПреподаватель имеет имя email должность.\nФакультет имеет название код.\nКафедра имеет название кабинет.\nКурс имеет название код кредиты.\nГруппа имеет номер год.\nЭкзамен имеет дату время кабинет.\nРезультат имеет оценку дату статус.\nФакультет содержит кафедры.\nКафедра содержит курсы.\nСтудент принадлежит группе.\nСтудент записывается на курсы.\nПреподаватель ведет курсы.\nСтудент сдает экзамены.\nЭкзамен связан с курсом.\nПреподаватель оценивает результаты.\nРезультат связан с экзаменом.",
    "control": "Клиент имеет имя email телефон.\nТовар имеет название артикул цену остаток.\nКатегория имеет название код.\nПоставщик имеет название телефон email.\nЗаказ имеет номер дату сумму статус.\nОплата имеет сумму дату способ статус.\nДоставка имеет адрес дату статус.\nКлиент оформляет заказы.\nЗаказ содержит товары.\nТовар относится к категории.\nПоставщик поставляет товары.\nЗаказ содержит оплаты.\nЗаказ содержит доставки."
};

let exampleDatabases = {
    library: "LibraryDemo",
    university: "UniversityProcess",
    control: "CommerceDemo",
};

button.addEventListener("click", analyze);
uploadFileButton.addEventListener("click", () => sourceFileInput.click());
sourceFileInput.addEventListener("change", loadSourceFile);
addEntityButton.addEventListener("click", addEntity);
addRelationButton.addEventListener("click", addRelation);
regenerateButton.addEventListener("click", regenerateSQL);
downloadSQLButton.addEventListener("click", downloadSQL);
downloadJSONButton.addEventListener("click", downloadJSON);
downloadMermaidButton.addEventListener("click", downloadMermaid);
downloadReportButton.addEventListener("click", downloadReport);
diagramZoomOutButton.addEventListener("click", () => setDiagramZoom(diagramZoom - 0.12));
diagramZoomResetButton.addEventListener("click", () => setDiagramZoom(1));
diagramZoomInButton.addEventListener("click", () => setDiagramZoom(diagramZoom + 0.12));
diagramFitButton.addEventListener("click", fitDiagramToViewport);
diagramViewport.addEventListener("wheel", zoomDiagramWithWheel, {passive: false});
databaseNameInput.addEventListener("input", handleModelInputChange);
document.querySelectorAll("[data-example]").forEach((exampleButton) => {
    exampleButton.addEventListener("click", () => {
        const key = exampleButton.dataset.example;
        textArea.value = examples[key] || textArea.value;
        databaseNameInput.value = exampleDatabases[key] || databaseNameInput.value;
        analyze();
    });
});
window.addEventListener("load", initializeExamples);

async function initializeExamples() {
    try {
        const response = await fetch("/domain-examples.json", {cache: "no-store"});
        if (response.ok) {
            const profiles = await response.json();
            examples = Object.fromEntries(Object.entries(profiles).map(([key, profile]) => [key, profile.text]));
            exampleDatabases = Object.fromEntries(Object.entries(profiles).map(([key, profile]) => [key, profile.database]));
        }
    } catch (error) {
        console.warn("Domain examples fallback is used:", error);
    }
    analyze();
}

async function loadSourceFile() {
    const file = sourceFileInput.files?.[0];
    if (!file) return;

    const extension = file.name.split(".").pop()?.toLowerCase() || "";
    if (!["txt", "json", "csv"].includes(extension)) {
        statusNode.textContent = "Неподдерживаемый формат файла. Загрузите .txt, .json или .csv.";
        sourceFileInput.value = "";
        return;
    }

    try {
        const content = await file.text();
        if (!content.trim()) {
            statusNode.textContent = "Загруженный файл пустой.";
            sourceFileInput.value = "";
            return;
        }

        textArea.value = content;
        databaseNameInput.value = databaseNameFromFile(file.name) || databaseNameInput.value || "database";
        statusNode.textContent = `Файл ${file.name} загружен. Выполняется анализ...`;
        await analyze();
    } catch (error) {
        statusNode.textContent = `Ошибка загрузки файла: ${error.message}`;
    } finally {
        sourceFileInput.value = "";
    }
}

function databaseNameFromFile(fileName) {
    const baseName = fileName.replace(/\.[^.]+$/, "").trim();
    return baseName.replace(/[_-]+/g, " ").replace(/\s+/g, " ").trim();
}

async function analyze() {
    const text = textArea.value.trim();
    if (!text) {
        statusNode.textContent = "Описание предметной области пустое.";
        return;
    }

    button.disabled = true;
    statusNode.textContent = "Выполняется семантический анализ...";

    try {
        const response = await fetch("/analyze", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify({text, database: {name: databaseNameInput.value.trim() || "database"}}),
        });

        if (!response.ok) throw new Error(await response.text());

        const result = await response.json();
        const diagnostics = result.diagnostics || [];
        currentModel = {
            database: result.database || {name: "database"},
            entities: result.entities || [],
            relations: result.relations || [],
        };
        currentAnalysis = {
            input: text,
            diagnostics,
            transformations: result.transformations || [],
            explanation: result.explanation || {candidates: []},
            sql: result.sql || "",
        };
        diagramPositions.clear();
        diagramNeedsCenter = true;
        modelDirty = false;
        renderDatabase(currentModel.database);
        renderEntities(currentModel.entities);
        renderRelations(currentModel.relations, currentModel.entities);
        renderDiagnostics(diagnostics);
        renderTrace(currentAnalysis.explanation?.candidates || []);
        renderTransformations(currentAnalysis.transformations);
        renderEvaluation(currentModel, currentAnalysis);
        renderSQL(currentAnalysis.sql);
        renderDiagram(currentModel.entities, currentModel.relations);
        statusNode.textContent = modelStatusText(currentModel, diagnostics);
    } catch (error) {
        statusNode.textContent = `Ошибка анализа: ${error.message}`;
    } finally {
        button.disabled = false;
    }
}

function renderDatabase(database) {
    databaseNameInput.value = database?.name || "database";
}

function renderEntities(entities) {
    entitiesNode.innerHTML = "";

    entities.forEach((entity, entityIndex) => {
        const card = document.createElement("article");
        card.className = "entity";
        card.innerHTML = `
            <label class="entity-name">
                <span>
                    Сущность
                    <button class="remove-entity-button" data-entity-index="${entityIndex}" data-action="remove-entity" type="button" title="Удалить сущность">x</button>
                </span>
                <input data-entity-index="${entityIndex}" data-field="name" value="${escapeAttribute(entity.name)}">
            </label>
        `;

        const attributes = entity.attributes || [];
        const attributesNode = document.createElement("div");
        attributesNode.className = "entity-attributes";
        if (!attributes.length) {
            const empty = document.createElement("div");
            empty.className = "attribute empty-attribute";
            empty.textContent = "Атрибуты не обнаружены";
            attributesNode.append(empty);
        }

        attributes.forEach((attr, attributeIndex) => {
            const row = document.createElement("div");
            row.className = "attribute";
            row.innerHTML = `
                <input data-entity-index="${entityIndex}" data-attribute-index="${attributeIndex}" data-field="attribute-name" value="${escapeAttribute(attr.name)}">
                <select data-entity-index="${entityIndex}" data-attribute-index="${attributeIndex}" data-field="attribute-type">
                    ${typeOptions(attr.type)}
                </select>
                <label class="attribute-flag" title="NOT NULL">
                    <input type="checkbox" data-entity-index="${entityIndex}" data-attribute-index="${attributeIndex}" data-field="attribute-required"${attr.required ? " checked" : ""}>
                    <span>Обяз.</span>
                </label>
                <label class="attribute-flag" title="UNIQUE">
                    <input type="checkbox" data-entity-index="${entityIndex}" data-attribute-index="${attributeIndex}" data-field="attribute-unique"${attr.unique ? " checked" : ""}>
                    <span>Уник.</span>
                </label>
                <button class="icon-button danger-button" data-entity-index="${entityIndex}" data-attribute-index="${attributeIndex}" data-action="remove-attribute" type="button" title="Удалить атрибут">x</button>
            `;
            attributesNode.append(row);
        });

        card.append(attributesNode);
        const addButton = document.createElement("button");
        addButton.className = "add-row-button";
        addButton.dataset.entityIndex = String(entityIndex);
        addButton.dataset.action = "add-attribute";
        addButton.type = "button";
        addButton.textContent = "+ Атрибут";
        card.append(addButton);
        entitiesNode.append(card);
    });

    entitiesNode.querySelectorAll("input, select").forEach((input) => {
        input.addEventListener("change", handleModelInputChange);
        input.addEventListener("input", handleModelInputChange);
    });
    entitiesNode.querySelectorAll("[data-action='add-attribute']").forEach((input) => input.addEventListener("click", addAttribute));
    entitiesNode.querySelectorAll("[data-action='remove-attribute']").forEach((input) => input.addEventListener("click", removeAttribute));
    entitiesNode.querySelectorAll("[data-action='remove-entity']").forEach((input) => input.addEventListener("click", removeEntity));
}

function renderRelations(relations, entities) {
    relationsNode.innerHTML = "";
    if (!relations.length) {
        relationsNode.textContent = "Связи не обнаружены";
        return;
    }

    const names = entities.map((entity) => entity.name);
    relations.forEach((relation, relationIndex) => {
        const row = document.createElement("article");
        row.className = "relation-row";
        row.innerHTML = `
            <select data-relation-index="${relationIndex}" data-field="relation-from">${entityOptions(names, relation.from)}</select>
            <select data-relation-index="${relationIndex}" data-field="relation-type">${relationTypeOptions(relation.type)}</select>
            <select data-relation-index="${relationIndex}" data-field="relation-to">${entityOptions(names, relation.to)}</select>
            <select data-relation-index="${relationIndex}" data-field="relation-cardinality">${cardinalityOptions(relation.cardinality)}</select>
            <button class="icon-button danger-button" data-relation-index="${relationIndex}" data-action="remove-relation" type="button" title="Удалить связь">x</button>
        `;
        relationsNode.append(row);
    });

    relationsNode.querySelectorAll("select").forEach((input) => input.addEventListener("change", handleModelInputChange));
    relationsNode.querySelectorAll("[data-action='remove-relation']").forEach((input) => input.addEventListener("click", removeRelation));
}

function renderSQL(sql) {
    sqlNode.textContent = sql || "-- SQL появится после анализа";
}

function renderDiagnostics(diagnostics) {
    diagnosticsNode.innerHTML = "";
    if (!diagnostics.length) {
        const item = document.createElement("div");
        item.className = "diagnostic diagnostic-info";
        item.textContent = "Проверка модели выполнена успешно.";
        diagnosticsNode.append(item);
        return;
    }

    for (const diagnostic of diagnostics) {
        const item = document.createElement("div");
        item.className = `diagnostic ${diagnosticClass(diagnostic.level)}`;
        item.textContent = diagnosticText(diagnostic);
        diagnosticsNode.append(item);
    }
}

function renderTrace(candidates) {
    traceNode.innerHTML = "";
    if (!candidates.length) {
        traceNode.textContent = "Срабатывания правил пока не найдены.";
        return;
    }

    for (const candidate of candidates) {
        const row = document.createElement("article");
        row.className = "trace-row";
        const target = candidate.target ? ` -> ${candidate.target}` : "";
        const owner = candidate.owner ? `${candidate.owner}.` : "";
        const confidence = Math.round((candidate.confidence || 0) * 100);
        row.innerHTML = `
            <div class="trace-main">
                <span class="trace-kind">${escapeHTML(candidate.kind)}</span>
                <strong>${escapeHTML(owner)}${escapeHTML(candidate.name)}${escapeHTML(target)}</strong>
                <span class="trace-confidence">${confidence}%</span>
            </div>
            <div class="trace-meta">
                <span>${escapeHTML(candidate.rule)}</span>
                <span>${escapeHTML(candidate.sourceText)}</span>
            </div>
        `;
        traceNode.append(row);
    }
}

function renderTransformations(transformations) {
    transformationsNode.innerHTML = "";
    if (!transformations.length) {
        transformationsNode.textContent = "Шаги преобразования пока не сформированы.";
        return;
    }

    for (const step of transformations) {
        const row = document.createElement("article");
        row.className = "transform-row";
        row.innerHTML = `
            <div class="transform-main">
                <span class="transform-stage">${escapeHTML(step.stage)}</span>
                <strong>${escapeHTML(step.source)} -> ${escapeHTML(step.target)}</strong>
            </div>
            <div class="transform-meta">
                <span>${escapeHTML(step.rule)}</span>
                <span>${escapeHTML(step.details || "")}</span>
            </div>
        `;
        transformationsNode.append(row);
    }
}

function renderEvaluation(model, analysis) {
    const entities = model.entities || [];
    const relations = model.relations || [];
    const diagnostics = analysis.diagnostics || [];
    const candidates = analysis.explanation?.candidates || [];
    const attributes = entities.reduce((sum, entity) => sum + (entity.attributes || []).length, 0);
    const errors = diagnostics.filter((diagnostic) => diagnostic.level === "error").length;
    const warnings = diagnostics.filter((diagnostic) => diagnostic.level === "warning").length;
    const info = diagnostics.filter((diagnostic) => diagnostic.level !== "error" && diagnostic.level !== "warning").length;
    const problemDiagnostics = diagnostics.filter((diagnostic) => diagnostic.level === "error" || diagnostic.level === "warning");
    const accepted = candidates.filter((candidate) => candidate.accepted !== false);
    const confidence = accepted.length
        ? Math.round(accepted.reduce((sum, candidate) => sum + (candidate.confidence || 0), 0) * 100 / accepted.length)
        : 0;
    const sqlReady = Boolean((analysis.sql || "").trim());
    const transformationCount = analysis.transformations?.length || 0;
    const modelScore = Math.max(0, Math.min(100,
        100
        - errors * 30
        - warnings * 8
        - (entities.length ? 0 : 24)
        - (relations.length ? 0 : 12)
        - (sqlReady ? 0 : 18),
    ));
    const qualityState = errors
        ? "Требуется исправление"
        : warnings
            ? "Есть предупреждения"
            : "Готово к использованию";
    const inputMode = accepted.length ? "Семантические правила" : "Структурированный ввод";

    evaluationNode.innerHTML = `
        <div class="metric-grid">
            ${metricCard("Сущности", entities.length, "Таблицы текущей модели")}
            ${metricCard("Атрибуты", attributes, "Поля текущей модели")}
            ${metricCard("Связи", relations.length, "Связи текущей ER-модели")}
            ${metricCard("Уверенность", accepted.length ? `${confidence}%` : "100%", accepted.length ? "Среднее по принятым правилам" : "Файл задан явно")}
            ${metricCard("Качество", `${modelScore}%`, qualityState)}
            ${metricCard("SQL", sqlReady ? "Готов" : "Нет", `${transformationCount} шагов преобразования`)}
        </div>
        <div class="current-evaluation">
            <div class="evaluation-summary">
                <strong>${escapeHTML(qualityState)}</strong>
                <span>${escapeHTML(inputMode)}. Ошибки: ${errors}, предупреждения: ${warnings}, информационные сообщения: ${info}.</span>
            </div>
            <div class="evaluation-checks">
                ${evaluationCheck("Сущности выделены", entities.length > 0)}
                ${evaluationCheck("Связи определены", relations.length > 0)}
                ${evaluationCheck("SQL-схема сформирована", sqlReady)}
                ${evaluationCheck("Критических ошибок нет", errors === 0)}
            </div>
            <div class="evaluation-issues">
                <strong>${problemDiagnostics.length ? "Что проверить" : "Сомнительных мест не найдено"}</strong>
                ${problemDiagnostics.length
                    ? problemDiagnostics.map((diagnostic) => evaluationIssue(diagnostic)).join("")
                    : `<span class="evaluation-issue is-info">Модель прошла проверку без ошибок и предупреждений.</span>`}
            </div>
        </div>
    `;
}

function metricCard(label, value, note) {
    return `
        <article class="metric-card">
            <span>${escapeHTML(label)}</span>
            <strong>${escapeHTML(value)}</strong>
            <small>${escapeHTML(note)}</small>
        </article>
    `;
}

function evaluationCheck(label, passed) {
    return `
        <span class="evaluation-check ${passed ? "is-passed" : "is-failed"}">
            ${passed ? "✓" : "!"} ${escapeHTML(label)}
        </span>
    `;
}

function evaluationIssue(diagnostic) {
    const isError = diagnostic.level === "error";
    return `
        <span class="evaluation-issue ${isError ? "is-error" : "is-warning"}">
            <b>${isError ? "Ошибка" : "Предупреждение"}</b>
            ${escapeHTML(diagnosticText(diagnostic))}
        </span>
    `;
}

async function regenerateSQL() {
    updateCurrentModelFromForm();
    regenerateButton.disabled = true;
    statusNode.textContent = "SQL пересобирается по отредактированной модели...";

    try {
        const response = await fetch("/generate-sql", {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(currentModel),
        });

        if (!response.ok) throw new Error(await response.text());

        const result = await response.json();
        const diagnostics = result.diagnostics || [];
        currentAnalysis = {
            ...currentAnalysis,
            diagnostics,
            transformations: result.transformations || [],
            sql: result.sql || "",
        };
        modelDirty = false;
        renderSQL(result.sql || "");
        renderDiagnostics(diagnostics);
        renderTransformations(currentAnalysis.transformations);
        renderEvaluation(currentModel, currentAnalysis);
        renderDiagram(currentModel.entities, currentModel.relations);
        renderRelations(currentModel.relations, currentModel.entities);
        statusNode.textContent = `SQL пересобран по отредактированной модели. ${diagnosticSummary(diagnostics)}`;
    } catch (error) {
        statusNode.textContent = `Ошибка генерации: ${error.message}`;
    } finally {
        regenerateButton.disabled = false;
    }
}

function handleModelInputChange() {
    updateCurrentModelFromForm();
    renderDiagram(currentModel.entities, currentModel.relations);
    renderEvaluation(currentModel, currentAnalysis);
    markModelDirty();
}

function markModelDirty() {
    modelDirty = true;
    statusNode.textContent = "Модель изменена. Пересоберите SQL, чтобы обновить DDL, диагностику и отчет преобразования.";
}

function updateCurrentModelFromForm() {
    const database = {...(currentModel.database || {}), name: databaseNameInput.value.trim() || "database"};
    const entities = JSON.parse(JSON.stringify(currentModel.entities || []));

    entitiesNode.querySelectorAll("[data-field='name']").forEach((input) => {
        entities[Number(input.dataset.entityIndex)].name = input.value.trim();
    });
    entitiesNode.querySelectorAll("[data-field='attribute-name']").forEach((input) => {
        entities[Number(input.dataset.entityIndex)].attributes[Number(input.dataset.attributeIndex)].name = input.value.trim();
    });
    entitiesNode.querySelectorAll("[data-field='attribute-type']").forEach((input) => {
        entities[Number(input.dataset.entityIndex)].attributes[Number(input.dataset.attributeIndex)].type = input.value;
    });
    entitiesNode.querySelectorAll("[data-field='attribute-required']").forEach((input) => {
        entities[Number(input.dataset.entityIndex)].attributes[Number(input.dataset.attributeIndex)].required = input.checked;
    });
    entitiesNode.querySelectorAll("[data-field='attribute-unique']").forEach((input) => {
        entities[Number(input.dataset.entityIndex)].attributes[Number(input.dataset.attributeIndex)].unique = input.checked;
    });

    const nameByOldName = new Map();
    (currentModel.entities || []).forEach((entity, index) => {
        const newName = entities[index].name;
        nameByOldName.set(entity.name, newName);
        if (entity.name !== newName && diagramPositions.has(entity.name)) {
            diagramPositions.set(newName, diagramPositions.get(entity.name));
            diagramPositions.delete(entity.name);
        }
    });

    const relations = (currentModel.relations || []).map((relation) => ({
        ...relation,
        from: nameByOldName.get(relation.from) || relation.from,
        to: nameByOldName.get(relation.to) || relation.to,
    }));

    relationsNode.querySelectorAll("[data-field='relation-from']").forEach((input) => {
        relations[Number(input.dataset.relationIndex)].from = nameByOldName.get(input.value) || input.value;
    });
    relationsNode.querySelectorAll("[data-field='relation-to']").forEach((input) => {
        relations[Number(input.dataset.relationIndex)].to = nameByOldName.get(input.value) || input.value;
    });
    relationsNode.querySelectorAll("[data-field='relation-type']").forEach((input) => {
        relations[Number(input.dataset.relationIndex)].type = input.value;
    });
    relationsNode.querySelectorAll("[data-field='relation-cardinality']").forEach((input) => {
        relations[Number(input.dataset.relationIndex)].cardinality = input.value;
    });

    currentModel = {database, entities, relations};
}

function addAttribute(event) {
    updateCurrentModelFromForm();
    currentModel.entities[Number(event.currentTarget.dataset.entityIndex)].attributes.push({name: "new_field", type: "TEXT", required: false, unique: false});
    renderEntities(currentModel.entities);
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function addEntity() {
    updateCurrentModelFromForm();
    const name = nextEntityName(currentModel.entities);
    currentModel.entities.push({name, attributes: []});
    renderEntities(currentModel.entities);
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function removeEntity(event) {
    updateCurrentModelFromForm();
    const index = Number(event.currentTarget.dataset.entityIndex);
    const entity = currentModel.entities[index];
    if (!entity) return;

    currentModel.entities.splice(index, 1);
    currentModel.relations = currentModel.relations.filter((relation) => relation.from !== entity.name && relation.to !== entity.name);
    renderEntities(currentModel.entities);
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function removeAttribute(event) {
    updateCurrentModelFromForm();
    currentModel.entities[Number(event.currentTarget.dataset.entityIndex)].attributes.splice(Number(event.currentTarget.dataset.attributeIndex), 1);
    renderEntities(currentModel.entities);
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function addRelation() {
    updateCurrentModelFromForm();
    if (currentModel.entities.length < 2) {
        statusNode.textContent = "Для добавления связи нужны минимум две сущности.";
        return;
    }
    currentModel.relations.push({from: currentModel.entities[0].name, to: currentModel.entities[1].name, type: "has", cardinality: "one-to-many"});
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function removeRelation(event) {
    updateCurrentModelFromForm();
    currentModel.relations.splice(Number(event.currentTarget.dataset.relationIndex), 1);
    renderRelations(currentModel.relations, currentModel.entities);
    renderDiagram(currentModel.entities, currentModel.relations);
    markModelDirty();
}

function nextEntityName(entities) {
    const names = new Set((entities || []).map((entity) => entity.name));
    let index = entities.length + 1;
    let name = `NewEntity${index}`;
    while (names.has(name)) {
        index += 1;
        name = `NewEntity${index}`;
    }
    return name;
}

function downloadSQL() {
    downloadText("database_schema.sql", sqlNode.textContent || "", "text/sql");
}

function downloadJSON() {
    updateCurrentModelFromForm();
    downloadText("database_model.json", JSON.stringify(currentModel, null, 2), "application/json");
}

function downloadMermaid() {
    updateCurrentModelFromForm();
    downloadText("er_model.mmd", buildMermaid(currentModel.entities, currentModel.relations), "text/plain");
}

function downloadReport() {
    updateCurrentModelFromForm();
    const report = {
        generatedAt: new Date().toISOString(),
        stale: modelDirty,
        input: currentAnalysis.input,
        model: currentModel,
        diagnostics: currentAnalysis.diagnostics,
        transformations: currentAnalysis.transformations,
        explanation: currentAnalysis.explanation,
        sql: sqlNode.textContent || currentAnalysis.sql || "",
    };
    downloadText("analysis_report.json", JSON.stringify(report, null, 2), "application/json");
}

function downloadText(filename, content, type) {
    const blob = new Blob([content], {type});
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = filename;
    document.body.append(link);
    link.click();
    link.remove();
    URL.revokeObjectURL(url);
}

function buildMermaid(entities, relations) {
    const lines = ["erDiagram"];
    for (const entity of entities) {
        lines.push(`    ${mermaidName(entity.name)} {`);
        lines.push("        int id PK");
        for (const attribute of entity.attributes || []) {
            lines.push(`        ${mermaidType(attribute.type)} ${mermaidName(attribute.name)}`);
        }
        lines.push("    }");
    }
    for (const relation of relations) {
        lines.push(`    ${mermaidName(relation.from)} ${mermaidCardinality(relation.cardinality)} ${mermaidName(relation.to)} : ${relation.type || "relates"}`);
    }
    return lines.join("\n");
}

function mermaidName(value) {
    return String(value || "item").trim().replaceAll(/[^\p{L}\p{N}_]+/gu, "_").replaceAll(/^_+|_+$/g, "").toUpperCase() || "ITEM";
}

function mermaidType(value) {
    const type = String(value || "TEXT").toUpperCase();
    if (type.includes("INT")) return "int";
    if (type.includes("NUMERIC")) return "decimal";
    if (type.includes("DATE")) return "date";
    if (type.includes("BOOL")) return "boolean";
    return "string";
}

function mermaidCardinality(value) {
    switch (value) {
        case "one-to-many":
            return "||--o{";
        case "many-to-one":
            return "}o--||";
        case "many-to-many":
            return "}o--o{";
        case "one-to-one":
            return "||--||";
        default:
            return "}o--o{";
    }
}

function typeOptions(selectedType) {
    const types = ["TEXT", "VARCHAR(255)", "VARCHAR(20)", "INTEGER", "NUMERIC(12,2)", "DATE", "TIME", "BOOLEAN"];
    return types.map((type) => `<option value="${escapeAttribute(type)}"${type === selectedType ? " selected" : ""}>${escapeHTML(type)}</option>`).join("");
}

function entityOptions(names, selectedName) {
    return names.map((name) => `<option value="${escapeAttribute(name)}"${name === selectedName ? " selected" : ""}>${escapeHTML(name)}</option>`).join("");
}

function relationTypeOptions(selectedType) {
    const types = ["has", "belongs_to", "contains", "associated_with"];
    return types.map((type) => `<option value="${escapeAttribute(type)}"${type === selectedType ? " selected" : ""}>${escapeHTML(type)}</option>`).join("");
}

function cardinalityOptions(selectedCardinality) {
    const cardinalities = ["one-to-one", "one-to-many", "many-to-one", "many-to-many", "unspecified"];
    return cardinalities.map((cardinality) => `<option value="${escapeAttribute(cardinality)}"${cardinality === selectedCardinality ? " selected" : ""}>${escapeHTML(cardinality)}</option>`).join("");
}

function diagnosticClass(level) {
    if (level === "error") return "diagnostic-error";
    if (level === "warning") return "diagnostic-warning";
    return "diagnostic-info";
}

function diagnosticText(diagnostic) {
    const message = String(diagnostic?.message || "");
    const translations = [
        [/^NLP analysis uses spaCy tokens, lemmas, POS tags and dependency metadata\.$/, "NLP-анализ выполнен через spaCy: использованы токены, леммы, части речи и зависимости."],
        [/^NLP service is unavailable or returned weak tags; local rule-based fallback is used\.$/, "NLP-сервис недоступен или вернул слабую разметку; использован локальный набор правил."],
        [/^Structured JSON input was parsed directly; NLP analysis was skipped\.$/, "Структурированный JSON разобран напрямую; NLP-анализ не выполнялся."],
        [/^Structured CSV input was parsed directly; NLP analysis was skipped\.$/, "Структурированный CSV разобран напрямую; NLP-анализ не выполнялся."],
        [/^Generated SQL passed internal DDL sanity checks\.$/, "Сгенерированный SQL прошел внутреннюю проверку DDL."],
        [/^Generated SQL is empty\.$/, "Сгенерированный SQL пустой."],
        [/^Database name is empty; SQL generation will skip CREATE DATABASE\.$/, "Имя базы данных не указано; команда CREATE DATABASE будет пропущена."],
        [/^No entities were detected; at least one table is required for a database schema\.$/, "Сущности не обнаружены; для схемы базы данных нужна хотя бы одна таблица."],
        [/^Entity with empty name was found\.$/, "Обнаружена сущность без имени."],
        [/^Several entities were detected, but no relations were found\. Add relation phrases or edit relations manually\.$/, "Обнаружено несколько сущностей, но связи не найдены. Добавьте фразы о связях или исправьте модель вручную."],
        [/^Relation with empty endpoint was found\.$/, "Обнаружена связь с пустым началом или концом."],
        [/^Generated SQL has a closing parenthesis without a matching opening parenthesis\.$/, "В SQL найдена закрывающая скобка без соответствующей открывающей."],
        [/^Generated SQL has unbalanced parentheses\.$/, "В SQL нарушен баланс скобок."],
    ];

    for (const [pattern, text] of translations) {
        if (pattern.test(message)) return text;
    }

    const replacements = [
        [/^Duplicate entity name "(.+)" was found\.$/, 'Повторяется имя сущности "$1".'],
        [/^Entities "(.+)" and "(.+)" produce the same SQL table name "(.+)"\.$/, 'Сущности "$1" и "$2" дают одинаковое имя SQL-таблицы "$3".'],
        [/^Entity "(.+)" uses a reserved SQL identifier; generator will rename it to "(.+)"\.$/, 'Сущность "$1" использует зарезервированное SQL-имя; генератор переименует ее в "$2".'],
        [/^Entity "(.+)" has no scalar attributes except the generated primary key\.$/, 'У сущности "$1" нет скалярных атрибутов, кроме автоматически созданного первичного ключа.'],
        [/^Entity "(.+)" has an attribute with empty name\.$/, 'У сущности "$1" есть атрибут без имени.'],
        [/^Entity "(.+)" has duplicate attribute "(.+)"\.$/, 'У сущности "$1" повторяется атрибут "$2".'],
        [/^Attributes "(.+)" and "(.+)" of entity "(.+)" produce the same SQL column name "(.+)"\.$/, 'Атрибуты "$1" и "$2" сущности "$3" дают одинаковое имя SQL-столбца "$4".'],
        [/^Entity "(.+)" has attribute "(.+)"; generated tables already contain primary key id\.$/, 'У сущности "$1" есть атрибут "$2"; таблицы уже содержат первичный ключ id.'],
        [/^Attribute "(.+)" of entity "(.+)" has unsupported type "(.+)"\.$/, 'Атрибут "$1" сущности "$2" имеет неподдерживаемый тип "$3".'],
        [/^Relation source "(.+)" does not match any entity\.$/, 'Источник связи "$1" не соответствует ни одной сущности.'],
        [/^Relation target "(.+)" does not match any entity\.$/, 'Цель связи "$1" не соответствует ни одной сущности.'],
        [/^Relation "(.+)" -> "(.+)" points to the same entity; check whether this is intentional\.$/, 'Связь "$1" -> "$2" указывает на ту же сущность. Проверьте, действительно ли это нужно.'],
        [/^Relation "(.+)" -> "(.+)" has unsupported cardinality "(.+)"\.$/, 'Связь "$1" -> "$2" имеет неподдерживаемую кардинальность "$3".'],
        [/^Generated SQL statement does not end with a semicolon: (.+)$/, 'SQL-выражение не заканчивается точкой с запятой: $1'],
        [/^Generated SQL contains duplicate CREATE TABLE for "(.+)"\.$/, 'SQL содержит повторный CREATE TABLE для "$1".'],
        [/^Foreign key references missing table "(.+)"\.$/, 'Внешний ключ ссылается на отсутствующую таблицу "$1".'],
        [/^Generated SQL contains duplicate CREATE INDEX "(.+)"\.$/, 'SQL содержит повторный CREATE INDEX "$1".'],
        [/^Index "(.+)" targets missing table "(.+)"\.$/, 'Индекс "$1" указывает на отсутствующую таблицу "$2".'],
        [/^Result was saved to PostgreSQL with storage key (.+)\.$/, 'Результат сохранен в PostgreSQL с ключом $1.'],
        [/^PostgreSQL result storage failed: (.+)$/, 'Не удалось сохранить результат в PostgreSQL: $1'],
    ];

    for (const [pattern, text] of replacements) {
        if (pattern.test(message)) return message.replace(pattern, text);
    }

    return message;
}

function diagnosticSummary(diagnostics) {
    const errors = diagnostics.filter((diagnostic) => diagnostic.level === "error").length;
    const warnings = diagnostics.filter((diagnostic) => diagnostic.level === "warning").length;
    if (errors) return `Ошибки: ${errors}, предупреждения: ${warnings}.`;
    if (warnings) return `Предупреждения: ${warnings}.`;
    return "Проверка пройдена.";
}

function modelStatusText(model, diagnostics) {
    return `База данных: ${model.database.name || "database"}. Сущности: ${model.entities.length}, связи: ${model.relations.length}. ${diagnosticSummary(diagnostics)}`;
}

function renderDiagram(entities, relations) {
    const viewportWidth = Math.max(diagramViewport?.clientWidth || 0, 720);
    diagramBaseSize = diagramCanvasSizeFor(entities.length, viewportWidth);
    const {width, height} = diagramBaseSize;
    diagram.setAttribute("viewBox", `0 0 ${width} ${height}`);
    applyDiagramZoom();
    diagram.innerHTML = "";

    ensureDiagramPositions(entities, relations, width, height);
    removeStaleDiagramPositions(entities);

    const relationLayer = svgGroup("diagram-relations");
    const nodeLayer = svgGroup("diagram-nodes");
    diagram.append(relationLayer, nodeLayer);

    const pairCounts = new Map();
    const pairIndexes = new Map();
    for (const relation of relations) {
        const key = relationPairKey(relation);
        pairCounts.set(key, (pairCounts.get(key) || 0) + 1);
    }

    for (const relation of relations) {
        const from = diagramPositions.get(relation.from);
        const to = diagramPositions.get(relation.to);
        if (!from || !to) continue;

        const key = relationPairKey(relation);
        const pairIndex = pairIndexes.get(key) || 0;
        pairIndexes.set(key, pairIndex + 1);
        const pairTotal = pairCounts.get(key) || 1;
        const offset = (pairIndex - (pairTotal - 1) / 2) * 34;

        const color = relationColor(relation.cardinality);
        relationEdge(relationLayer, from, to, relation, color, offset);
    }

    for (const entity of entities) {
        const position = diagramPositions.get(entity.name);
        entityNode(nodeLayer, position.x, position.y, entity);
    }

    if (diagramNeedsCenter) {
        requestAnimationFrame(centerDiagramViewport);
        diagramNeedsCenter = false;
    }
}

function ensureDiagramPositions(entities, relations, width, height) {
    let addedPosition = false;
    entities.forEach((entity, index) => {
        if (diagramPositions.has(entity.name)) {
            const position = diagramPositions.get(entity.name);
            diagramPositions.set(entity.name, clampPosition(position, width, height));
            return;
        }

        diagramPositions.set(entity.name, gridDiagramPosition(index, entities.length, width, height));
        addedPosition = true;
    });

    if (addedPosition && relations.length) {
        arrangeDiagramByRelations(entities, relations, width, height);
    }
}

function diagramCanvasSizeFor(entityCount, viewportWidth) {
    const columns = Math.max(3, Math.ceil(Math.sqrt(Math.max(entityCount, 1))));
    const rows = Math.ceil(Math.max(entityCount, 1) / columns);
    return {
        width: Math.max(1800, Math.round(viewportWidth * 2.25), columns * 360 + 760),
        height: Math.max(1100, rows * 300 + 620),
    };
}

function diagramColumnCount(entityCount, width) {
    if (entityCount <= 1) return 1;
    const maxColumns = Math.max(1, Math.floor((width - 120) / 250));
    return Math.min(maxColumns, Math.ceil(Math.sqrt(entityCount)));
}

function gridDiagramPosition(index, total, width, height) {
    const columns = diagramColumnCount(total, width);
    const rows = Math.ceil(total / columns);
    const row = Math.floor(index / columns);
    const col = index % columns;
    const rowSize = row === rows - 1 ? total - row * columns : columns;
    const centeredCol = col + (columns - rowSize) / 2;
    const horizontalGap = (width - diagramComponentMargin.x * 2) / Math.max(columns - 1, 1);
    const verticalGap = (height - diagramComponentMargin.y * 2) / Math.max(rows - 1, 1);

    return clampPosition({
        x: columns === 1 ? width / 2 : diagramComponentMargin.x + centeredCol * horizontalGap,
        y: rows === 1 ? height / 2 : diagramComponentMargin.y + row * verticalGap,
    }, width, height);
}

function arrangeDiagramByRelations(entities, relations, width, height) {
    const entityNames = new Set(entities.map((entity) => entity.name));
    const adjacency = new Map(entities.map((entity) => [entity.name, new Set()]));
    for (const relation of relations) {
        if (!entityNames.has(relation.from) || !entityNames.has(relation.to)) continue;
        adjacency.get(relation.from).add(relation.to);
        adjacency.get(relation.to).add(relation.from);
    }

    const components = connectedComponents(entities, adjacency)
        .sort((left, right) => right.length - left.length);
    const componentSlots = componentCenters(components.length, width, height);

    components.forEach((component, index) => {
        arrangeComponent(component, adjacency, componentSlots[index], width, height);
    });
}

function connectedComponents(entities, adjacency) {
    const visited = new Set();
    const components = [];
    for (const entity of entities) {
        if (visited.has(entity.name)) continue;
        const queue = [entity.name];
        const component = [];
        visited.add(entity.name);
        while (queue.length) {
            const name = queue.shift();
            component.push(name);
            for (const next of adjacency.get(name) || []) {
                if (visited.has(next)) continue;
                visited.add(next);
                queue.push(next);
            }
        }
        components.push(component);
    }
    return components;
}

function componentCenters(count, width, height) {
    if (count <= 1) {
        return [{x: width / 2, y: height / 2}];
    }
    return Array.from({length: count}, (_, index) => gridDiagramPosition(index, count, width, height));
}

function arrangeComponent(component, adjacency, center, width, height) {
    if (component.length === 1) {
        diagramPositions.set(component[0], clampPosition(center, width, height));
        return;
    }

    const sorted = [...component].sort((left, right) => {
        const degreeDelta = (adjacency.get(right)?.size || 0) - (adjacency.get(left)?.size || 0);
        return degreeDelta || left.localeCompare(right);
    });
    const hub = sorted[0];
    const directNeighbors = sorted.filter((name) => adjacency.get(hub)?.has(name));
    const others = sorted.filter((name) => name !== hub && !adjacency.get(hub)?.has(name));
    diagramPositions.set(hub, clampPosition(center, width, height));

    const firstRing = directNeighbors.slice(0, 8);
    placeRing(firstRing, center, 360, width, height, -Math.PI / 2);

    if (directNeighbors.length > 8) {
        placeRing(directNeighbors.slice(8), center, 600, width, height, -Math.PI / 2 + Math.PI / 8);
    }

    if (others.length) {
        placeSecondaryNodes(others, adjacency, center, width, height);
    }
}

function placeRing(names, center, radius, width, height, startAngle) {
    names.forEach((name, index) => {
        const angle = startAngle + (Math.PI * 2 * index) / Math.max(names.length, 1);
        diagramPositions.set(name, clampPosition({
            x: center.x + Math.cos(angle) * radius,
            y: center.y + Math.sin(angle) * radius,
        }, width, height));
    });
}

function placeSecondaryNodes(names, adjacency, center, width, height) {
    const groupedByAnchor = new Map();
    for (const name of names) {
        const anchor = [...(adjacency.get(name) || [])].find((candidate) => diagramPositions.has(candidate)) || null;
        if (!groupedByAnchor.has(anchor)) groupedByAnchor.set(anchor, []);
        groupedByAnchor.get(anchor).push(name);
    }

    for (const [anchor, group] of groupedByAnchor.entries()) {
        const anchorPosition = anchor ? diagramPositions.get(anchor) : center;
        const angle = Math.atan2(anchorPosition.y - center.y, anchorPosition.x - center.x);
        const direction = {
            x: Math.cos(angle || -Math.PI / 2),
            y: Math.sin(angle || -Math.PI / 2),
        };
        const tangent = {x: -direction.y, y: direction.x};
        const spread = 270;

        group.forEach((name, index) => {
            const offset = (index - (group.length - 1) / 2) * spread;
            diagramPositions.set(name, clampPosition({
                x: anchorPosition.x + direction.x * 280 + tangent.x * offset,
                y: anchorPosition.y + direction.y * 280 + tangent.y * offset,
            }, width, height));
        });
    }
}

function removeStaleDiagramPositions(entities) {
    const names = new Set(entities.map((entity) => entity.name));
    for (const name of diagramPositions.keys()) {
        if (!names.has(name)) {
            diagramPositions.delete(name);
        }
    }
}

function clampPosition(position, width, height) {
    return {
        x: Math.min(Math.max(position.x, diagramNodeMargin.x), width - diagramNodeMargin.x),
        y: Math.min(Math.max(position.y, diagramNodeMargin.y), height - diagramNodeMargin.y),
    };
}

function relationPairKey(relation) {
    return [relation.from, relation.to].sort().join("::");
}

function relationEdge(layer, from, to, relation, color, offset) {
    const dx = to.x - from.x;
    const dy = to.y - from.y;
    const distance = Math.hypot(dx, dy) || 1;
    const normalX = -dy / distance;
    const normalY = dx / distance;
    const start = {
        x: from.x + (dx / distance) * 122 + normalX * offset,
        y: from.y + (dy / distance) * 62 + normalY * offset,
    };
    const end = {
        x: to.x - (dx / distance) * 122 + normalX * offset,
        y: to.y - (dy / distance) * 62 + normalY * offset,
    };
    const label = {
        x: (start.x + end.x) / 2 + normalX * 22,
        y: (start.y + end.y) / 2 + normalY * 22,
    };

    line(layer, start.x, start.y, end.x, end.y, color);
    labelText(layer, label.x, label.y, relation.cardinality || relation.type, color);
}

function entityNode(layer, centerX, centerY, entity) {
    const attrs = entity.attributes || [];
    const shown = attrs.slice(0, 4);
    const width = 238;
    const height = 76 + shown.length * 20 + (attrs.length > shown.length ? 20 : 0);
    const x = centerX - width / 2;
    const y = centerY - height / 2;

    const group = svgGroup("diagram-entity");
    group.dataset.entityName = entity.name;
    group.setAttribute("tabindex", "0");
    group.setAttribute("role", "button");
    group.setAttribute("aria-label", `Move ${entity.name}`);
    group.addEventListener("pointerdown", startDiagramDrag);

    rect(group, x, y, width, height, "#ffffff", "#2563eb");
    rect(group, x, y, width, 46, "#eaf4ff", "#2563eb");
    rect(group, x, y, 7, height, "#2563eb", "#2563eb");
    text(group, x + 18, y + 27, entity.name, "15px", "#172033", "800", "start");
    text(group, x + width - 18, y + 27, "TABLE", "10px", "#2563eb", "900", "end");
    text(group, x + 20, y + 62, "id", "12px", "#2563eb", "800", "start");
    text(group, x + width - 20, y + 62, "PK", "11px", "#0f766e", "900", "end");

    shown.forEach((attribute, index) => {
        const flags = `${attribute.required ? "!" : ""}${attribute.unique ? " U" : ""}`;
        const rowY = y + 84 + index * 20;
        text(group, x + 20, rowY, attribute.name, "12px", "#334155", "700", "start");
        text(group, x + width - 20, rowY, `${mermaidType(attribute.type)}${flags}`, "11px", "#64748b", "800", "end");
    });
    if (attrs.length > shown.length) {
        text(group, x + 20, y + 84 + shown.length * 20, `+${attrs.length - shown.length} еще`, "11px", "#64748b", "800", "start");
    }

    layer.append(group);
}

function relationColor(cardinality) {
    switch (cardinality) {
        case "many-to-many":
            return "#7c3aed";
        case "one-to-many":
        case "many-to-one":
            return "#0f766e";
        case "one-to-one":
            return "#2563eb";
        default:
            return "#64748b";
    }
}

function startDiagramDrag(event) {
    if (event.button !== 0) return;
    const entityName = event.currentTarget.dataset.entityName;
    const position = diagramPositions.get(entityName);
    if (!position) return;

    const pointer = diagramPointer(event);
    activeDrag = {
        entityName,
        offsetX: pointer.x - position.x,
        offsetY: pointer.y - position.y,
    };
    event.currentTarget.setPointerCapture(event.pointerId);
    diagram.classList.add("is-dragging");
    event.stopPropagation();
    event.preventDefault();
}

diagram.addEventListener("pointermove", (event) => {
    if (!activeDrag) return;

    const pointer = diagramPointer(event);
    diagramPositions.set(activeDrag.entityName, clampPosition({
        x: pointer.x - activeDrag.offsetX,
        y: pointer.y - activeDrag.offsetY,
    }, diagramBaseSize.width, diagramBaseSize.height));
    renderDiagram(currentModel.entities || [], currentModel.relations || []);
});

diagram.addEventListener("pointerup", stopDiagramDrag);
diagram.addEventListener("pointercancel", stopDiagramDrag);
diagram.addEventListener("pointerleave", stopDiagramDrag);
diagramViewport.addEventListener("pointerdown", startDiagramPan);
diagramViewport.addEventListener("pointermove", moveDiagramPan);
diagramViewport.addEventListener("pointerup", stopDiagramPan);
diagramViewport.addEventListener("pointercancel", stopDiagramPan);
diagramViewport.addEventListener("pointerleave", stopDiagramPan);

function stopDiagramDrag() {
    if (!activeDrag) return;
    activeDrag = null;
    diagram.classList.remove("is-dragging");
}

function diagramPointer(event) {
    const point = diagram.createSVGPoint();
    point.x = event.clientX;
    point.y = event.clientY;
    return point.matrixTransform(diagram.getScreenCTM().inverse());
}

function startDiagramPan(event) {
    if (event.button !== 0 || event.target.closest(".diagram-entity")) return;
    activePan = {
        x: event.clientX,
        y: event.clientY,
        scrollLeft: diagramViewport.scrollLeft,
        scrollTop: diagramViewport.scrollTop,
    };
    diagramViewport.setPointerCapture(event.pointerId);
    diagramViewport.classList.add("is-panning");
    event.preventDefault();
}

function moveDiagramPan(event) {
    if (!activePan) return;
    diagramViewport.scrollLeft = activePan.scrollLeft - (event.clientX - activePan.x);
    diagramViewport.scrollTop = activePan.scrollTop - (event.clientY - activePan.y);
}

function stopDiagramPan() {
    if (!activePan) return;
    activePan = null;
    diagramViewport.classList.remove("is-panning");
}

function setDiagramZoom(value, focalEvent = null) {
    const previousZoom = diagramZoom;
    const nextZoom = Math.min(Math.max(value, 0.38), 2.4);
    if (Math.abs(nextZoom - previousZoom) < 0.001) return;

    let anchor = null;
    if (focalEvent) {
        const viewportRect = diagramViewport.getBoundingClientRect();
        anchor = {
            x: (diagramViewport.scrollLeft + focalEvent.clientX - viewportRect.left) / previousZoom,
            y: (diagramViewport.scrollTop + focalEvent.clientY - viewportRect.top) / previousZoom,
            offsetX: focalEvent.clientX - viewportRect.left,
            offsetY: focalEvent.clientY - viewportRect.top,
        };
    }

    diagramZoom = nextZoom;
    applyDiagramZoom();

    if (anchor) {
        diagramViewport.scrollLeft = anchor.x * diagramZoom - anchor.offsetX;
        diagramViewport.scrollTop = anchor.y * diagramZoom - anchor.offsetY;
    }
}

function zoomDiagramWithWheel(event) {
    if (event.target.closest("select, input, textarea, button")) return;
    event.preventDefault();
    const direction = event.deltaY > 0 ? -1 : 1;
    const multiplier = direction > 0 ? 1.11 : 0.9;
    setDiagramZoom(diagramZoom * multiplier, event);
}

function fitDiagramToViewport() {
    const viewportWidth = Math.max(diagramViewport?.clientWidth || 0, 1);
    const viewportHeight = Math.max(diagramViewport?.clientHeight || 0, 1);
    const bounds = diagramContentBounds();
    if (!bounds) {
        setDiagramZoom(1);
        centerDiagramViewport();
        return;
    }

    const scaleX = (viewportWidth - 72) / bounds.width;
    const scaleY = (viewportHeight - 72) / bounds.height;
    const nextZoom = Math.max(0.72, Math.min(scaleX, scaleY, 1.18));
    setDiagramZoom(nextZoom);
    requestAnimationFrame(() => centerDiagramViewport(bounds));
}

function applyDiagramZoom() {
    diagram.style.width = `${Math.round(diagramBaseSize.width * diagramZoom)}px`;
    diagram.style.height = `${Math.round(diagramBaseSize.height * diagramZoom)}px`;
    diagramZoomResetButton.textContent = `${Math.round(diagramZoom * 100)}%`;
}

function centerDiagramViewport(bounds = null) {
    const target = bounds || diagramContentBounds();
    if (!target) {
        diagramViewport.scrollLeft = Math.max(0, (diagram.scrollWidth - diagramViewport.clientWidth) / 2);
        diagramViewport.scrollTop = Math.max(0, (diagram.scrollHeight - diagramViewport.clientHeight) / 2);
        return;
    }

    const contentCenterX = (target.minX + target.width / 2) * diagramZoom;
    const contentCenterY = (target.minY + target.height / 2) * diagramZoom;
    diagramViewport.scrollLeft = Math.max(0, contentCenterX - diagramViewport.clientWidth / 2);
    diagramViewport.scrollTop = Math.max(0, contentCenterY - diagramViewport.clientHeight / 2);
}

function diagramContentBounds() {
    if (!diagramPositions.size) return null;

    let minX = Infinity;
    let minY = Infinity;
    let maxX = -Infinity;
    let maxY = -Infinity;
    for (const position of diagramPositions.values()) {
        minX = Math.min(minX, position.x - 170);
        minY = Math.min(minY, position.y - 130);
        maxX = Math.max(maxX, position.x + 170);
        maxY = Math.max(maxY, position.y + 140);
    }

    return {
        minX: Math.max(0, minX),
        minY: Math.max(0, minY),
        width: Math.max(300, Math.min(diagramBaseSize.width, maxX) - Math.max(0, minX)),
        height: Math.max(240, Math.min(diagramBaseSize.height, maxY) - Math.max(0, minY)),
    };
}

function svgGroup(className) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "g");
    if (className) node.setAttribute("class", className);
    return node;
}

function rect(parent, x, y, width, height, fill, stroke) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "rect");
    node.setAttribute("x", x);
    node.setAttribute("y", y);
    node.setAttribute("width", width);
    node.setAttribute("height", height);
    node.setAttribute("rx", "8");
    node.setAttribute("fill", fill);
    node.setAttribute("stroke", stroke);
    node.setAttribute("stroke-width", "1.5");
    parent.append(node);
}

function line(parent, x1, y1, x2, y2, stroke) {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "line");
    node.setAttribute("x1", x1);
    node.setAttribute("y1", y1);
    node.setAttribute("x2", x2);
    node.setAttribute("y2", y2);
    node.setAttribute("stroke", stroke);
    node.setAttribute("stroke-width", "2.2");
    node.setAttribute("stroke-linecap", "round");
    parent.append(node);
}

function labelText(parent, x, y, value, fill) {
    const group = svgGroup("diagram-label");
    const labelWidth = Math.max(72, String(value || "").length * 7.2 + 18);
    rect(group, x - labelWidth / 2, y - 13, labelWidth, 26, "#ffffff", "#dbe4ef");
    text(group, x, y + 4, value, "12px", fill, "800");
    parent.append(group);
}

function text(parent, x, y, value, size, fill, weight = "500", anchor = "middle") {
    const node = document.createElementNS("http://www.w3.org/2000/svg", "text");
    node.setAttribute("x", x);
    node.setAttribute("y", y);
    node.setAttribute("text-anchor", anchor);
    node.setAttribute("font-size", size);
    node.setAttribute("font-weight", weight);
    node.setAttribute("fill", fill);
    node.textContent = value;
    parent.append(node);
    return node;
}

function escapeHTML(value) {
    return String(value).replaceAll("&", "&amp;").replaceAll("<", "&lt;").replaceAll(">", "&gt;").replaceAll('"', "&quot;");
}

function escapeAttribute(value) {
    return escapeHTML(value).replaceAll("'", "&#039;");
}
