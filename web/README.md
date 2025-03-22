# web

A new Flutter project.

## Getting Started

This project is a starting point for a Flutter application.

A few resources to get you started if this is your first Flutter project:

- [Lab: Write your first Flutter app](https://docs.flutter.dev/get-started/codelab)
- [Cookbook: Useful Flutter samples](https://docs.flutter.dev/cookbook)

For help getting started with Flutter development, view the
[online documentation](https://docs.flutter.dev/), which offers tutorials,
samples, guidance on mobile development, and a full API reference.

## How to contribute

Make a changes in UI. Build a project for web and copy files to static directory using following command

```
 rm -rf configo/static/*
 cp -r web/build/web/* configo/static/
```

so that it will be available in go server at endpoint http://localhost:8085/configo/web