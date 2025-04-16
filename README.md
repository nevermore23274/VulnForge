# VulnForge CLI

VulnForge CLI is a lightweight Go-based command-line tool that lets you scrape and transform cybersecurity intelligence sources into datasets formatted for instruction-tuned LLMs.

Currently, it supports scraping from Abuse.ch's ThreatFox API and converting indicators of compromise (IOCs) into JSON or CSV datasets for ML applications.

---

## 🚀 Features

- Scrape IOCs from ThreatFox (more sources coming soon)
- Save raw data as `.jsonl`
- Convert to instruction-tuning datasets (JSON or CSV)
- Supports `--dry-run` for previewing conversion output
- Built-in config system for setting defaults (e.g., tag, format, limit)
- Versioning support
- Compatible with GPU or CPU training environments
- Clean Docker-based development and training workflows

---

## ⚡ Quick Start

Just want to see it work? Here's the fastest path to generating a dataset and training a model.

### 1. Build both environments
```bash
docker compose build --build-arg UNSLOTH_DEVICE=cpu vulnforge trainer
```

### 2. Run VulnForge to scrape and convert
```bash
docker compose run --rm vulnforge go run ./main.go scrape --tag honeypot
```

```bash
# Replace with actual date or glob wildcard
FILENAME=$(ls output/threatfox_iocs_*.jsonl | head -n1)
```

```bash
docker compose run --rm vulnforge go run ./main.go convert \
  --infile $FILENAME \
  --outfile output/quick_dataset.json \
  --limit 20
```

### 3. Train a model using the trainer container
```bash
docker compose run --rm trainer python scripts/train_lora.py   --model gpt2   --dataset output/quick_dataset.json   --output_dir output/lora-gpt2   --device cpu
```

That's it! You’ve created a dataset and fine-tuned a model.

---

## 📦 Installation

### 🛠️ VulnForge CLI Environment (dataset tooling)

```bash
docker compose build vulnforge
```

Then start it:
```bash
docker compose run --rm vulnforge
```

### 🧠 Model Training Environment (LoRA fine-tuning)

#### 🖥️ CPU-only setup:
```bash
docker compose build --build-arg UNSLOTH_DEVICE=cpu trainer
```

### 🚀 GPU-enabled setup:
```bash
docker compose build --build-arg UNSLOTH_DEVICE=cuda trainer
```

Then start the container:
```bash
docker compose run --rm trainer
```

---

## 🧪 Usage

To view available commands and flags at any time:

```bash
./vulnforge --help
./vulnforge <command> --help
```

---

## 📁 Dataset Operations (VulnForge CLI)

### 📥 Scrape from ThreatFox

```bash
./vulnforge scrape --tag honeypot
```

Generates: `output/threatfox_iocs_<date>_honeypot.jsonl`

Or use default config:

```bash
./vulnforge config set default_tag honeypot
./vulnforge scrape
```

### 🔁 Convert to LLM Dataset

```bash
./vulnforge convert   --infile output/threatfox_iocs_<date>_honeypot.jsonl   --outfile output/instruction_dataset_<date>.json   --format json --limit 50
```

Or preview with:

```bash
./vulnforge convert --infile output/*.jsonl --dry-run
```

---

## 🧠 Training Models (inside `trainer` container)

Run the training script from the shared `/workspace` directory:

```bash
python scripts/train_lora.py   --model gpt2   --dataset output/instruction_dataset_<date>.json   --output_dir output/lora-gpt2   --device cpu
```

Set `--device cuda` to use GPU (if enabled).

---

## ⚙️ Configuration

### 🔧 Set Defaults

```bash
./vulnforge config set default_format csv
./vulnforge config set default_tag honeypot
./vulnforge config set default_limit 50
```

### 📋 View or Reset Config

```bash
./vulnforge config get
./vulnforge config reset
```

---

## 📌 Version Info

```bash
./vulnforge version
```

---

## 📖 License

Apache 2.0

---

## 💬 Credits

- Abuse.ch ThreatFox API
- Cobra CLI Framework
- Hugging Face Transformers / PEFT
- Docker Compose + cross-platform support
