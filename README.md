# Calendar App Backend

A backend service made with [Fiber](https://gofiber.io/) to manage accounts, jwt authentications, events crud, and serving a react app in the public folder.

Note: If you would like to change the front end app, just upload the build of your application in ``` ./public ``` folder.

## Instalation

1. Clone this project
2. Copy .env.template file and rename it to .env
3. Create a mongo database and get the connection url
4. Change env mongo database url for the one you just created
5. Run the project with command:
```
go run .
```

## Frontend
This backend runs the build version of this [front-end project](https://github.com/Sebas3270/calendar-app), if you would like to change something in the ui or ux, use the link to go the complete react app, there you can change everything you want, just don't forget to run the build again and change in public folder.