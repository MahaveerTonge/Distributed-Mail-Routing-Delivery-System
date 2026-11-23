# Distributed Mail Routing & Delivery System

A horizontally scalable, fault-tolerant distributed mail routing pipeline built in **Go**, **Python**, and **Terraform** on **AWS**. The system decouples high-throughput message ingestion from downstream delivery using distributed queues, idempotent state tracking, circuit-broken AI payload classification, and end-to-end observability.

## Architecture Overview

```text
Clients / Load Gen
       │
       ▼
Application Load Balancer (ALB)
       │
       ▼
Stateless Go Edge Ingestion Nodes ──► DynamoDB (Idempotent State Store)
       │
       ▼
Amazon SQS (Main Queue + Dead-Letter Queue)
       ├──► Go Delivery Worker Pool ──► Mock SMTP Sink (Configurable Delay/Bounces)
       └──► Python AI Classifier Pool (Tier-1 Linear + Tier-2 Quantized Model + Circuit Breaker)
       │
       ▼
Observability Stack (OpenTelemetry + Prometheus + Grafana)
