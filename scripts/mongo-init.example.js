/* init.js – будет выполнен только при первом запуске контейнера MongoDB */
const ROOT_USER = "your_mongo_user" // From .env
const ROOT_PWD = "your_mongo_password" // From .env
const APP_DB = "your_mongo_db" // From .env

db = db.getSiblingDB(APP_DB)

db.createUser({
	user: ROOT_USER,
	pwd: ROOT_PWD,
	roles: [{ role: "readWrite", db: APP_DB }]
})

db.createCollection("requests")
db.createCollection("chapters")

db.requests.createIndex({ title: 1 }, { unique: false })
