/**
 AthenaeumModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class AthenaeumModel: ObservableObject {
    @Published private(set) var books: Set<Book>
    @Published private(set) var selectedBookId: Book.ID?

    init() {
        books = []
    }
}

extension AthenaeumModel {
    func selectBook(_ book: Book) {
        selectBook(id: book.id)
    }

    func selectBook(id: Book.ID) {
        selectedBookId = id
    }

    func addBook(from _: URL) {}
}
