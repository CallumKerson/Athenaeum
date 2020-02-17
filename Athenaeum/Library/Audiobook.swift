/**
 Audiobook.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class Audiobook {
    var title: String
    var author: String
    var file: URL
    var narrator: String?
    var publicationDate: String?
    var isbn: String?
    var summary: String?
    var entry: String?
    var series: Series?

    init(title: String,
         author: String,
         file: URL,
         narrator: String? = nil,
         publicationDate: String? = nil,
         isbn: String? = nil,
         summary: String? = nil,
         entry: String? = nil,
         series: Series? = nil) {
        self.title = title
        self.author = author
        self.file = file
        self.narrator = narrator
        self.publicationDate = publicationDate
        self.isbn = isbn
        self.summary = summary
        self.entry = entry
        self.series = series
    }
}

extension Audiobook: Identifiable {
//    let id = UUID()
}

extension Audiobook: CustomStringConvertible {
    var description: String {
        "Audiobook (\(title) by \(author))"
    }
}

extension Audiobook: Equatable {
    static func == (lhs: Audiobook, rhs: Audiobook) -> Bool {
        guard lhs.title == rhs.title else {
            return false
        }
        guard lhs.author == rhs.author else {
            return false
        }
        return true
    }
}

extension Audiobook: Hashable {
    func hash(into hasher: inout Hasher) {
        hasher.combine(title)
        hasher.combine(author)
    }
}
