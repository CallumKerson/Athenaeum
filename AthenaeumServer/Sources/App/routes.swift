/**
 routes.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Fluent
import Vapor

func routes(_ app: Application) throws {
    app.get { req in
        req.view.render("index", ["title": "Hello Vapor!"])
    }

    app.get("hello") { _ -> String in
        "Hello, world!"
    }

    app.get("feed.rss") { _ -> String in
        generateFeed()
    }

    try app.register(collection: TodoController())
}
