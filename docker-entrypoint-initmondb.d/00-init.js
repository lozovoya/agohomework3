db.createUser({
    user: "app",
    pwd: "pass",
    roles: [
        {
            role: "readWrite",
            db: "db",
        }
    ]
});


db.createCollection('users', {
    validator: {
        $jsonSchema: {
            bsonType: 'object',
            required:['userid', 'name', 'account', 'suggestions'],
            properties: {
                userid: {
                    bsonType: 'int'
                },
                name: {
                    bsonType: 'string'
                },
                account: {
                    bsonType: 'string'
                },
                suggestions: {
                    bsonType: 'array',
                    items: {
                        bsonType: 'object',
                        required: ['sugid', 'icon', 'title', 'link'],
                        properties: {
                            sugid: {
                                bsonType: 'int'
                            },
                            icon: {
                                bsonType: 'string'
                            },
                            title: {
                                bsonType: 'string'
                            },
                            link: {
                                bsonType: 'string'
                            }
                        }
                    }
                },
                operations: {
                    bsonType: 'array',
                    items: {
                        bsonType: 'object',
                        required: ['oppid', 'amount', 'description'],
                        properties: {
                            oppid: {
                                bsonType: 'int'
                            },
                            amount: {
                                bsonType: 'int'
                            },
                            description: {
                                bsonType: 'string'
                            }
                        }
                    }
                }
            }
        }
    }
    }
);



