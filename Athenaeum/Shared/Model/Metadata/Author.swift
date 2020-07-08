/**
 Author.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

/// A struct describing a book author
struct Author: Equatable, Codable, Hashable, ExpressibleByStringLiteral {
    public let firstNames: String?
    public let lastName: String

    public var fullName: String {
        if let firstName = firstNames {
            return "\(firstName) \(lastName)"
        } else {
            return lastName
        }
    }

    public init(firstNames: String? = nil, lastName: String) {
        self.firstNames = firstNames
        self.lastName = lastName
    }

//    /// Initalises Author from string, including special cases
//    /// - Parameter fullName: Full name of author
//    public init(fullName: String) {
//        // Exceptions
//        if fullName == "Ursula K. Le Guin" {
//            self.init(firstNames: "Ursula K.", lastName: "Le Guin")
//        } else {
//            self.init(stringLiteral: fullName)
//        }
//    }

    /// Initalises Author from string
    /// - Parameter stringLiteral: full name of author
    public init(stringLiteral value: String) {
        if value == "Ursula K. Le Guin" {
            self.init(firstNames: "Ursula K.", lastName: "Le Guin")
        } else {
            var names = value.components(separatedBy: " ")
            if let lastName = names.last {
                names.removeLast(1)
                self.init(firstNames: names.joined(separator: " "), lastName: lastName)
            } else {
                self.init(firstNames: nil, lastName: value)
            }
        }
    }
}

extension Array where Element == Author {
    /// Gets all authors in a single string
    var author: String {
        map { $0.fullName }.joined(separator: ", ")
            .replacingLastOccurrenceOfString(",", with: " &")
    }
}
