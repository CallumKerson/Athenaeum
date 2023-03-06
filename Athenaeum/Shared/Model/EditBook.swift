/**
 EditBook.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

struct EditBook: Encodable, Identifiable {
    var id: String?
    var title: String?
    var author: [Person]?
    var summary: String?
    var releaseDate: Date?
    @NullEncodable var series: Series?
    var genre: [Genre]?

    init(from book: Book) {
        self.id = book.id
        self.title = book.title
        self.author = book.author
        self.summary = book.summary
        self.releaseDate = book.releaseDate
        self.series = book.series
        self.genre = book.genre
    }

    init(
        id: String? = nil,
        title: String? = nil,
        author: [Person]? = nil,
        summary: String? = nil,
        releaseDate: Date? = nil,
        series: Series? = nil,
        genre: [Genre]? = nil
    ) {
        self.id = id
        self.title = title
        self.author = author
        self.summary = summary
        self.releaseDate = releaseDate
        self.series = series
        self.genre = genre
    }
}

@propertyWrapper
struct NullEncodable<T>: Encodable where T: Encodable {
    var wrappedValue: T?

    init(wrappedValue: T?) {
        self.wrappedValue = wrappedValue
    }

    func encode(to encoder: Encoder) throws {
        var container = encoder.singleValueContainer()
        switch self.wrappedValue {
        case let .some(value): try container.encode(value)
        case .none: try container.encodeNil()
        }
    }
}
