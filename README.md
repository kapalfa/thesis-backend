# A web framework for collaborative code development
This repository consists of the backend of my thesis "A web framework for collaborative code development"

## What the application does

The application creates a collaborative coding environment for developers. Every user who has an account in the app, can create projects and store their files there. 
The app offers a code editor for editing and viewing purposes, a live chat between collaborators, GitHub integration and many more features.

## Why did I use Go

- Performance: Go is a compiled language
- Concurrency: Go has built-in support for concurrency through goroutines and channels
- Easier deployment: Go produces statically compiled binaries, so you don't need to install a runtime environment on the target system

## Challenges and Future Developments

A terminal was implemented for our application so that the users could compile and run their code. However when we deployed to Cloud Storage, the terminal didn't have direct access to the users' projects
to compile and run them. Our goal is to implement a fully functional terminal that will complete the development experience.

## Contents
### Controllers
- Authentication: we used access and refresh tokens and jwt authentication
- Files: functions for creating, deleting, reading, uploading and saving files and folders
- Projects: functions for creating, deleting, searching, copying projects, etc
- Github Controllers: functions for initializing, downloading, commiting changes to a GitHub remote repository, etc
- Chat: implementation of live chat between the collaborators of the same project
### Models
- User
- Project
- Access : connects the user with the project table
