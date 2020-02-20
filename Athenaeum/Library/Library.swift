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
    let shouldUpdateAudiobooks = PassthroughSubject<Void, Never>()
    let shouldUpdatePublisher: AnyPublisher<Void, Never>
    
//    private let concurrentLibraryQueue =
//        DispatchQueue(
//            label: "com.umbra.Athenaeum.Library",
//            attributes: .concurrent)
    
    var 🎧📚: [Audiobook] = [] {
        willSet {
        self.objectWillChange.send()
        }
    }
    
    private var didChangeCancellable: AnyCancellable?
    
    init() {
        shouldUpdatePublisher = shouldUpdateAudiobooks.eraseToAnyPublisher()
        log.info("Initialising Library")
        self.🎧📚 = Library.getRepository().getAll()
        didChangeCancellable = shouldUpdateAudiobooks.sink(receiveValue: {_ in
            log.info("Updating current books in library")
            self.🎧📚 = Library.getRepository().getAll()
        } )
    }
    
    public func shelve(book: Audiobook) {
//        concurrentLibraryQueue.async(flags: .barrier) { [weak self] in
//            guard self != nil else {
//                return
//            }
//            
            let repository = Library.getRepository();
            do {
                try repository.insert(item: book)
            } catch {
                log.error("Cannot add \(book) to Library")
                log.debug(error)
            }
//        }
        DispatchQueue.main.async {
            self.shouldUpdateAudiobooks.send()
        }
    }
    
    private static func getRepository() -> RealmRepository<Audiobook> {
        RealmRepository()
    }
}
