/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Combine
import Foundation

final class RepositoryLibrary: Library {
    static var global = RepositoryLibrary()

    let objectWillChange = ObservableObjectPublisher()

    var 🎧📚: [Audiobook] = [] {
        willSet {
            objectWillChange.send()
        }
    }

    private var didChangeCancellable: AnyCancellable?

    init() {
        log.info("Initialising Library")
        🎧📚 = RepositoryLibrary.getRepository().getAll()
    }

    public func shelve(book: Audiobook) {
        if let bookFile = book as? AudiobookFile {
            let repository = RepositoryLibrary.getRepository()
            do {
                try repository.insert(item: bookFile)
            } catch {
                log.error("Cannot add \(book) to Library")
                log.debug(error)
            }
            DispatchQueue.main.async {
                self.🎧📚 = RepositoryLibrary.getRepository().getAll()
                self.objectWillChange.send()
            }
        } else {
            log.error("Only AudiobookFiles supported in RepositoryLibrary")
        }
    }

    private static func getRepository() -> RealmRepository<AudiobookFile> {
        RealmRepository()
    }
}
