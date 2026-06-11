from __future__ import annotations

import math
import textwrap
from pathlib import Path

from PIL import Image, ImageDraw, ImageFont


ROOT = Path(__file__).resolve().parent
OUT = ROOT / "architecture_scheme.png"

W, H = 2200, 1400

FONT = Path(r"C:\Windows\Fonts\arial.ttf")
BOLD = Path(r"C:\Windows\Fonts\arialbd.ttf")


def font(size: int, bold: bool = False) -> ImageFont.FreeTypeFont:
    return ImageFont.truetype(str(BOLD if bold else FONT), size)


F_TITLE = font(30, True)
F_H1 = font(24, True)
F_H2 = font(19, True)
F_TEXT = font(17)
F_SMALL = font(15)
F_TINY = font(13)


img = Image.new("RGB", (W, H), "#f5f7fb")
draw = ImageDraw.Draw(img)


def rounded(box, fill, outline="#cbd5e1", width=2, radius=18):
    draw.rounded_rectangle(box, radius=radius, fill=fill, outline=outline, width=width)


def text(x, y, value, fnt=F_TEXT, fill="#172033", anchor="la"):
    draw.text((x, y), value, font=fnt, fill=fill, anchor=anchor)


def centered(x, y, value, fnt=F_TEXT, fill="#172033"):
    draw.text((x, y), value, font=fnt, fill=fill, anchor="mm")


def wrap_lines(value: str, fnt, max_width: int) -> list[str]:
    words = value.split()
    lines: list[str] = []
    current = ""
    for word in words:
        candidate = word if not current else current + " " + word
        bbox = draw.textbbox((0, 0), candidate, font=fnt)
        if bbox[2] - bbox[0] <= max_width:
            current = candidate
        else:
            if current:
                lines.append(current)
            current = word
    if current:
        lines.append(current)
    return lines


def paragraph(x, y, value, fnt=F_SMALL, fill="#334155", max_width=250, line_gap=5):
    cursor = y
    for line in wrap_lines(value, fnt, max_width):
        text(x, cursor, line, fnt, fill)
        cursor += fnt.size + line_gap
    return cursor


def pill(x, y, w, h, label, fill, fg="#ffffff"):
    draw.rounded_rectangle((x, y, x + w, y + h), radius=12, fill=fill)
    centered(x + w / 2, y + h / 2, label, F_TINY, fg)


def arrow(start, end, color="#111827", width=3, dashed=False):
    x1, y1 = start
    x2, y2 = end
    if dashed:
        dx, dy = x2 - x1, y2 - y1
        length = math.hypot(dx, dy)
        if length == 0:
            return
        ux, uy = dx / length, dy / length
        pos = 0
        while pos < length - 18:
            a = pos
            b = min(pos + 16, length - 18)
            draw.line((x1 + ux * a, y1 + uy * a, x1 + ux * b, y1 + uy * b), fill=color, width=width)
            pos += 28
    else:
        draw.line((x1, y1, x2, y2), fill=color, width=width)

    angle = math.atan2(y2 - y1, x2 - x1)
    size = 14
    p1 = (x2, y2)
    p2 = (x2 - size * math.cos(angle - 0.42), y2 - size * math.sin(angle - 0.42))
    p3 = (x2 - size * math.cos(angle + 0.42), y2 - size * math.sin(angle + 0.42))
    draw.polygon([p1, p2, p3], fill=color)


def card(box, title, body, tag=None, fill="#ffffff", outline="#cbd5e1", tag_color="#2563eb"):
    x1, y1, x2, y2 = box
    rounded(box, fill, outline, 2, 18)
    if tag:
        pill(x1 + 18, y1 + 16, 110, 26, tag, tag_color)
        title_y = y1 + 56
    else:
        title_y = y1 + 24
    text(x1 + 20, title_y, title, F_H2, "#111827")
    paragraph(x1 + 20, title_y + 34, body, F_SMALL, "#475569", max_width=(x2 - x1 - 40))


# Header
text(70, 48, "Архитектура программной подсистемы автоматизированного проектирования", F_TITLE, "#0f172a")
text(70, 86, "структуры базы данных", F_TITLE, "#0f172a")
text(72, 130, "Локальное клиент-серверное приложение: веб-интерфейс, Go-бэкенд, Python NLP-сервис, PostgreSQL и экспорт артефактов.", F_TEXT, "#475569")

# Legend
rounded((1620, 40, 2110, 132), "#ffffff", "#d7dee9", 2, 16)
text(1644, 64, "Обозначения", F_H2, "#111827")
arrow((1648, 104), (1728, 104), "#111827", 3)
text(1746, 95, "REST / JSON", F_SMALL, "#334155")
arrow((1888, 104), (1968, 104), "#64748b", 3, dashed=True)
text(1986, 95, "опционально", F_SMALL, "#334155")

# Main zones
rounded((70, 180, 610, 640), "#eef6ff", "#8aa7d6", 3, 22)
pill(96, 204, 150, 30, "Вход и клиент", "#2563eb")
text(96, 258, "Пользовательский слой", F_H1, "#0f172a")
paragraph(96, 298, "Пользователь вводит описание предметной области или загружает файл. Интерфейс отображает ER-модель, SQL-код, диагностику и оценку результата.", F_TEXT, "#334155", 430)

card((110, 470, 290, 590), "Текст", "Описание предметной области на естественном языке.", fill="#ffffff", outline="#b7c9e8")
card((330, 470, 510, 590), "Файлы", "TXT, JSON и CSV для ручного ввода и интеграции.", fill="#ffffff", outline="#b7c9e8")

rounded((700, 180, 1400, 1050), "#fff7ed", "#d19a6a", 3, 24)
pill(730, 204, 155, 30, "Go backend", "#b86b35")
text(730, 258, "Серверная часть и pipeline проектирования", F_H1, "#0f172a")
paragraph(730, 298, "Go-приложение принимает запросы, координирует анализ, формирует модель базы данных, выполняет проверку и генерирует итоговые артефакты.", F_TEXT, "#334155", 590)

pipeline = [
    ("1", "HTTP-обработчики", "Раздача web-интерфейса, /analyze, /generate-sql, обмен JSON."),
    ("2", "Определение типа входа", "Текстовое описание, TXT, структурированный JSON или CSV."),
    ("3", "Семантический анализ", "Нормализация, сущности, атрибуты, связи и кардинальность."),
    ("4", "Валидация модели", "Ошибки, предупреждения и сомнительные места текущего результата."),
    ("5", "Генерация артефактов", "SQL DDL, JSON-модель, Mermaid-диаграмма и отчет анализа."),
]

y = 405
for num, title, body in pipeline:
    rounded((765, y, 1335, y + 96), "#fff1cf" if int(num) % 2 else "#fff8df", "#d4a76d", 2, 14)
    draw.ellipse((790, y + 28, 832, y + 70), fill="#b86b35")
    centered(811, y + 49, num, F_H2, "#ffffff")
    text(860, y + 20, title, F_H2, "#111827")
    paragraph(860, y + 52, body, F_SMALL, "#475569", 430)
    if num != "5":
        arrow((1050, y + 96), (1050, y + 128), "#9a5b2e", 3)
    y += 128

rounded((1500, 190, 2110, 475), "#ecfdf5", "#75a77d", 3, 22)
pill(1530, 218, 145, 30, "Python service", "#15803d")
text(1530, 274, "NLP-сервис", F_H1, "#0f172a")
paragraph(1530, 316, "FastAPI + spaCy выполняют токенизацию, лемматизацию, определение частей речи и зависимостей. Результат возвращается Go-сервису для дальнейшего анализа.", F_TEXT, "#334155", 520)
text(1530, 428, "Модели: ru_core_news_sm, en_core_web_sm", F_SMALL, "#166534")

rounded((1500, 540, 2110, 815), "#f0fdf4", "#75a77d", 3, 22)
pill(1530, 568, 145, 30, "PostgreSQL", "#5f8066")
text(1530, 624, "Хранение результатов", F_H1, "#0f172a")
paragraph(1530, 666, "Сохраняются исходное описание, модель, сущности, атрибуты, связи, SQL-код, диагностика и шаги преобразования. Хранилище используется для фиксации выполненного анализа.", F_TEXT, "#334155", 520)

rounded((1500, 890, 2110, 1145), "#eef6ff", "#7c96c4", 3, 22)
pill(1530, 918, 145, 30, "Artifacts", "#536f9a")
text(1530, 974, "Выходные артефакты", F_H1, "#0f172a")
paragraph(1530, 1016, "Итоговые результаты доступны в интерфейсе и для экспорта: SQL DDL, JSON-модель, Mermaid ER-диаграмма и отчет анализа.", F_TEXT, "#334155", 520)

# UI node
rounded((96, 360, 560, 430), "#ffffff", "#b7c9e8", 2, 16)
text(120, 386, "Веб-интерфейс: HTML / CSS / JavaScript", F_H2, "#111827")
text(120, 414, "экран проектирования, правка модели, экспорт", F_SMALL, "#475569")

# Arrows
arrow((610, 395), (700, 395), "#111827", 3)
text(630, 368, "POST /analyze", F_TINY, "#334155")

arrow((1400, 520), (1500, 350), "#64748b", 3, dashed=True)
text(1422, 420, "NLP-запрос", F_TINY, "#475569")

arrow((1400, 812), (1500, 682), "#111827", 3)
text(1414, 744, "SaveAnalysis", F_TINY, "#334155")

arrow((1400, 930), (1500, 1018), "#111827", 3)
text(1415, 976, "SQL / JSON / отчет", F_TINY, "#334155")

arrow((1500, 1104), (610, 1104), "#111827", 3)
arrow((610, 1104), (610, 430), "#111827", 3)
text(910, 1078, "результат возвращается в интерфейс", F_TINY, "#334155")

# Bottom note
rounded((70, 1200, 1400, 1330), "#ffffff", "#d7dee9", 2, 18)
text(96, 1234, "Логика разделения слоев", F_H2, "#111827")
paragraph(
    96,
    1268,
    "Интерфейс отвечает за взаимодействие с пользователем; Go-бэкенд выполняет проектную процедуру; Python-сервис используется для NLP-разметки; PostgreSQL фиксирует результат; экспортные артефакты применяются для дальнейшей демонстрации и использования.",
    F_TEXT,
    "#334155",
    1210,
)

img.save(OUT, quality=95)
print(OUT)
