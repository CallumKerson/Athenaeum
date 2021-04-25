/**
 Book.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Foundation

struct Series: Codable {
    let entry: Double?
    let title: String?
}

struct Person: Codable, Hashable {
    let givenNames: String
    let familyName: String
}

extension Person: CustomStringConvertible {
    var description: String {
        "\(self.givenNames) \(self.familyName)"
    }
}

struct Book: Codable, Identifiable {
    let id: String
    let title: String
    let author: [Person]
    let summary: String?
    let releaseDate: String?
    let series: Series?
}

extension Book {
    var authorString: String {
        if self.author.count == 1 {
            return self.author[0].description
        } else {
            let andAuthor = self.author[self.author.count - 1]
            var commaSeparatedAuthorStrings = [String]()
            for person in self.author.dropLast() {
                commaSeparatedAuthorStrings.append(person.description)
            }
            return "\(commaSeparatedAuthorStrings.joined(separator: ", ")) and \(andAuthor)"
        }
    }
}
