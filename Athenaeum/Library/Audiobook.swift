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
    var series: (title: String, entry: String)?

    init(title: String,
         author: String,
         file: URL,
         narrator: String? = nil,
         publicationDate: String? = nil,
         isbn: String? = nil,
         summary: String? = nil,
         series: (title: String, entry: String)? = nil) {
        self.title = title
        self.author = author
        self.file = file
        self.narrator = narrator
        self.publicationDate = publicationDate
        self.isbn = isbn
        self.summary = summary
        self.series = series
    }
}

extension Audiobook: Identifiable {
//    let id = UUID()
}

extension Audiobook: CustomStringConvertible {
    var description: String {
        "\(title) by \(author)"
    }
}
