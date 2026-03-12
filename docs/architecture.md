# SentinelX Architecture

SentinelX is designed as a modular security platform.

Architecture flow:

Traffic Sources
↓
Traffic Collectors
↓
Event Pipeline
↓
Detection Engines
↓
Threat Intelligence
↓
Response Engine
↓
Event Storage
↓
Dashboards

Modules:

collector
Collects network traffic and telemetry.

pipeline
Transfers events across the system.

detection
Analyzes events and detects attacks.

response
Triggers automated defense actions.

storage
Stores events and logs.

ui
Provides dashboards for visualization.