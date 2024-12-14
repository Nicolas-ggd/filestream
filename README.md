# FileStream: Event-Driven File Upload Engine
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream?ref=badge_shield)


FileStream is an open-source project, which is aim to gain more experience in golang, the purpose of this project is to build robust file upload systems that are efficient, resumable, and ready for real-time integrations.

## Purpose of this project

This project was born from a desire to learn Go deeply while building something practical and useful. My aim is to enhance my skills and gain real world experience by working on an open-source project that can also attract contributions from others who share the same passion. This project is learning journey for me and developers also which decide to collaborating and create a reusable engine.

## Roadmap

Here are some exciting features in the pipeline:
- Chunked uploads - Upload large files in smaller, manageable chunks.
- WebSocket Notifications - Real time updates for uploading process.
- Event driven - Publish events via NATS or other event system.
- Storage - In starting level it's good to support Minio for example.
- Customizable options - Customizable options is like remove metadata or using virus scan.
- Resumable upload - Resume interrupted uploads without starting over.

## How to contribute

Here’s how you can get involved:
1. Report issues: If you find any bug or issue, please open a [issue](https://github.com/Nicolas-ggd/filestream/issues)
2. Fork and code: Check out the [open issues](https://github.com/Nicolas-ggd/filestream/pulls) and submit a pull request.

## Project setup

1. Clone the repository:
    ```
    git clone https://github.com/Nicolas-ggd/filestream
    ```
2. Install dependencies:
    ```
   make dependencies
   ```
3. Run tests:
    ```
   make test
   ```
4. To start application run:
   ```
   make start
   ```
   
## License
FileStream is open-source software licensed under the MIT License.

[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FNicolas-ggd%2Ffilestream?ref=badge_large)