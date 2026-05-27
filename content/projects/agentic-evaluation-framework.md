---
title: "Agentic Evaluation Framework"
description: "A framework for evaluating how well LLMs perform open-ended tasks in simulated environments."
date: 2026-05-15
external_url: "https://github.com/juliusbroomfield"
---

Evaluating Large Language Models (LLMs) on static, multiple-choice benchmarks is no longer sufficient to gauge their capability in real-world applications. This project introduces a robust, extensible framework for orchestrating and evaluating agentic behavior in interactive, multi-step environments.

### Key Capabilities

*   **Multi-step environments**: Simulate bash terminals, web browsers, and file structures.
*   **State-based tracking**: Track environment state mutations rather than just text generation output.
*   **Safety validation**: Enforce run-time safety policies and catch malicious commands before execution.

---

### Framework Architecture

```
┌─────────────────┐       executes       ┌───────────────┐
│  Agent (LLM)    ├─────────────────────>│  Environment  │
└────────┬────────┘                      └───────┬───────┘
         │                                       │
         │ monitors behavior                     │ returns observations
         ▼                                       ▼
┌─────────────────┐                      ┌───────────────┐
│ Safety Monitor  │                      │ State Tracker │
└─────────────────┘                      └───────────────┘
```

---

### Baseline Performance Metrics

Here is a performance comparison of popular LLM backends evaluated across our agent tasks:

| Model Backend | Coding Success (%) | Planning Horizon (Steps) | Re-planning Accuracy (%) |
| :--- | :---: | :---: | :---: |
| **GPT-4o** | 74.2% | 15 | 88.5% |
| **Gemini 1.5 Pro** | 76.5% | 18 | 91.2% |
| **Claude 3.5 Sonnet** | 82.0% | 22 | 93.4% |
| **Llama 3 (70B)** | 61.8% | 10 | 79.0% |

---

### Getting Started

To initialize the evaluation benchmark locally:

```python
from agent_eval import EvaluationSuite, BashEnvironment

# Define environment parameters
env = BashEnvironment(sandbox_path="./sandbox", timeout=30)

# Instantiate the suite with a list of target tasks
suite = EvaluationSuite(
    tasks=["bash_manipulation", "file_edit"],
    environment=env
)

# Run evaluation on agent policy
metrics = suite.evaluate(agent_fn=my_llm_agent_policy)
print(f"Evaluation Completed! Success Rate: {metrics.success_rate:.2%}")
```
