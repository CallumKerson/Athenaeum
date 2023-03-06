/**
 Endpoint.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

struct Endpoint {
    var path: String
    var queryItems: [URLQueryItem]?
}

extension Endpoint {
    var athenaeumURL: URL? {
        guard var components = URLComponents(string: UserDefaults.standard
            .string(forKey: HOST_DEFAULTS_KEY) ?? "")
        else {
            return nil
        }
        components.path = "/api/v1" + self.path
        if let queryItems = queryItems {
            components.queryItems = queryItems
        }

        guard let url = components.url else {
            preconditionFailure("Invalid URL components: \(components)")
        }

        return url
    }

    var headers: [String: Any] {
        let username = UserDefaults.standard.string(forKey: USERNAME_DEFAULTS_KEY) ?? ""
        let password = UserDefaults.standard.string(forKey: PASSWORD_DEFAULTS_KEY) ?? ""
        return [
            "app-id": "Athenaeum",
            "Authorization": "Basic \(Data("\(username):\(password)".utf8).base64EncodedString())",
        ]
    }
}

extension Endpoint {
    static var books: Self {
        Endpoint(path: "/books/")
    }

    static func book(withID id: String) -> Self {
        Endpoint(path: "/books/\(id)")
    }

    static func genre(_ genre: Genre) -> Self {
        Endpoint(path: "/genres/\(genre.rawValue)")
    }
}
