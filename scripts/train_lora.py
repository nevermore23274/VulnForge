import argparse
import json
import torch
from transformers import AutoTokenizer, AutoModelForCausalLM, TrainingArguments, Trainer, DataCollatorForLanguageModeling
from datasets import Dataset
from peft import get_peft_model, LoraConfig, TaskType
from accelerate import Accelerator

# Parse CLI arguments
parser = argparse.ArgumentParser()
parser.add_argument("--model", type=str, default="gpt2", help="Base model to fine-tune")
parser.add_argument("--dataset", type=str, required=True, help="Path to VulnForge-formatted JSON dataset")
parser.add_argument("--output_dir", type=str, default="./lora-model", help="Directory to save model")
parser.add_argument("--device", type=str, default="auto", help="Device to use: 'cpu', 'cuda', or 'auto'")
args = parser.parse_args()

# Load dataset
with open(args.dataset, "r") as f:
    raw_data = json.load(f)

hf_dataset = Dataset.from_list(raw_data)

# Tokenizer
tokenizer = AutoTokenizer.from_pretrained(args.model)
if tokenizer.pad_token is None:
    tokenizer.pad_token = tokenizer.eos_token

def tokenize(example):
    prompt = f"{example['instruction']}\n{example['input']}\n{example['output']}"
    return tokenizer(
        prompt,
        padding="max_length",
        truncation=True,
        max_length=512
    )

tokenized_dataset = hf_dataset.map(tokenize, remove_columns=hf_dataset.column_names)

# Load model
model = AutoModelForCausalLM.from_pretrained(args.model)

# Apply LoRA
peft_config = LoraConfig(
    task_type=TaskType.CAUSAL_LM,
    inference_mode=False,
    r=8,
    lora_alpha=32,
    lora_dropout=0.1
)
model = get_peft_model(model, peft_config)

# Training arguments
training_args = TrainingArguments(
    output_dir=args.output_dir,
    per_device_train_batch_size=1,
    num_train_epochs=1,
    save_total_limit=1,
    logging_steps=10,
    save_steps=50,
    report_to="none",
    #evaluation_strategy="no"
)

# Accelerator handles device assignment
accelerator = Accelerator()
device = args.device if args.device in ["cpu", "cuda"] else accelerator.device
model.to(device)

# Trainer setup
trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=tokenized_dataset,
    tokenizer=tokenizer,
    data_collator=DataCollatorForLanguageModeling(tokenizer=tokenizer, mlm=False),
)

# Train
trainer.train()
trainer.save_model(args.output_dir)
print(f"✅ Model saved to {args.output_dir}")