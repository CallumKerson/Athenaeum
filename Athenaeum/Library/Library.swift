/**
 Library.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import Combine

final class Library: ObservableObject {
    
    static var global = Library()
    
    let objectWillChange = ObservableObjectPublisher()
    
    var 🎧📚: [Audiobook] = [] {
        willSet {
        self.objectWillChange.send()
        }
    }
    
    private var didChangeCancellable: AnyCancellable?
    
    init() {
        log.info("Initialising Library")
        self.🎧📚 = Library.getRepository().getAll()
    }
    
    public func shelve(book: Audiobook) {
            let repository = Library.getRepository();
            do {
                try repository.insert(item: book)
            } catch {
                log.error("Cannot add \(book) to Library")
                log.debug(error)
            }
        DispatchQueue.main.async {
            self.🎧📚 = Library.getRepository().getAll()
            self.objectWillChange.send()
        }
    }
    
    private static func getRepository() -> RealmRepository<Audiobook> {
        RealmRepository()
    }
}
