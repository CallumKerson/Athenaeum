/**
 AthenaeumModel.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

class AthenaeumModel: ObservableObject {
    @Published private(set) var books: [Book]
    @Published private(set) var selectedBookId: Book.ID?

    init() {
        let url = FileManager.default.urls(for: .musicDirectory, in: .userDomainMask)[0]
            .appendingPathComponent("MotOE.m4b")

        var newBook = Book(metadata: BookMetadata.fromAudiobook(audiobook: url)!)
        newBook.audio = url
        self.books = [newBook]
    }
}

extension AthenaeumModel {
    func selectBook(_ book: Book) {
        self.selectBook(id: book.id)
    }

    func selectBook(id: Book.ID) {
        self.selectedBookId = id
    }

    func addAudiobook(from url: URL) {
        guard url.pathExtension == "m4b" else {
            return
        }
        guard let metadata = BookMetadata.fromAudiobook(audiobook: url) else { return }
        var newBook = Book(metadata: metadata)
        newBook.audio = url
        self.books.append(newBook)
        print("Added new book \(newBook.metadata.title)")
    }
}
