db.createUser (
    {
        user: "carscrap",
        pwd: "mysuperpassword",
        roles: [
            {
                role: "readWrite",
                db: "carscrap"
            }
        ]
    }
)
