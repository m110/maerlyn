# Maerlyn

This project aims to provide a simple monitoring service for remote servers.

Featuring:

* Single-binary drop-in solution.
* Exposing server's health over a REST interface.
* Lightweight and fast.
* Configured by flags and env. vars instead of config files.

Checks could include:

* CPU usage (obviously)
* Memory usage
* System load
* Disk space usage
* I/O and network throughput.

The application can be run in two modes:

* Serving as HTTP server, binding on chosen port.
* Display all checks and quit (self-test for diagnostic purposes).
