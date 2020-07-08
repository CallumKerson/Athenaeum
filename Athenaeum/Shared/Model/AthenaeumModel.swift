/**
 AthenaeumModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation

class AthenaeumModel: ObservableObject {
    @Published private(set) var books: [Book]
    @Published private(set) var selectedBookId: Book.ID?

    init() {
        self.books = []
    }
}

extension AthenaeumModel {
    func selectBook(_ book: Book) {
        selectBook(id: book.id)
    }

    func selectBook(id: Book.ID) {
        selectedBookId = id
    }

    func addAudiobook(from url: URL) {
        guard url.pathExtension == "m4b" else {
            return
        }
        guard let metadata = BookMetadata.fromAudiobook(audiobook: url) else { return }
        var newBook = Book(metadata: metadata)
        newBook.audio = url
        books.append(newBook)
    }
}
