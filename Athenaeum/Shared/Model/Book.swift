/**
 Book.swift
 Copyright (c) 2021 Callum Kerr-Edwards
 */

import Foundation

struct Series: Codable {
    let entry: Double
    let title: String
}

extension Series: CustomStringConvertible {
    var description: String {
        if self.entry.truncatingRemainder(dividingBy: 1.0) == 0.0 {
            return "\(self.title) Book \(Int(self.entry))"
        } else {
            return "\(self.title) Book \(self.entry)"
        }
    }
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
    let releaseDate: Date?
    let genre: [Genre]?
    let series: Series?

    init(
        id: String,
        title: String,
        author: [Person],
        summary: String? = nil,
        releaseDate: Date? = nil,
        genre: [Genre]? = nil,
        series: Series? = nil
    ) {
        self.id = id
        self.title = title
        self.author = author
        self.summary = summary
        self.releaseDate = releaseDate
        self.genre = genre
        self.series = series
    }
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

extension Book {
    var shortReleaseDate: String? {
        let formatter = DateFormatter()
        formatter.dateStyle = .short
        guard let releaseDate = self.releaseDate else {
            logger
                .error(
                    "Cannot format shortReleaseDate, the releaseDate of book with \(self.id) is nil"
                )
            return nil
        }
        return formatter.string(from: releaseDate)
    }

    var mediumReleaseDate: String? {
        let formatter = DateFormatter()
        formatter.dateStyle = .medium
        guard let releaseDate = self.releaseDate else {
            logger
                .error(
                    "Cannot format mediumReleaseDate, the releaseDate of book with \(self.id) is nil"
                )
            return nil
        }
        return formatter.string(from: releaseDate)
    }
}

struct Books: Codable {
    var books: [Book]
}
