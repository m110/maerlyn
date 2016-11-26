# Maerlyn

This project aims to provide a simple monitoring service for remote servers. Now, that's of course a common task, already solved by many applications. Maerlyn is thought to be as easy as possible to deploy and do one thing really well. It exposes crucial health information of the device that can be accessed by another service (sending alerts or anything else).

## Overview

* Single-binary drop-in solution.
* Exposing server's health over a REST interface.
* Lightweight and fast.
* Configured by flags and env. vars instead of config files (see first bullet).

Checks could include:

* CPU usage (obviously)
* Memory usage
* System load
* Disk space usage
* I/O and network throughput.

The application can be run in two modes:

* Serving as HTTP server, binding on chosen port.
* Display all checks and quit (self-test for diagnostic purposes).
