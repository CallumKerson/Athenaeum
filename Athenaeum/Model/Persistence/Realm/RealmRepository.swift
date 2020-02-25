/**
 RealmRepository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation
import RealmSwift

class RealmRepository<RepositoryObject>: Repository
    where RepositoryObject: Entity,
    RepositoryObject.StoreType: Object {
    typealias RealmObject = RepositoryObject.StoreType

    private let subject = PassthroughSubject<DatabaseAction, Never>()
    let publisher: AnyPublisher<DatabaseAction, Never>

    private let realm: Realm

    init() {
        self.realm = try! Realm()
        self.publisher = self.subject.eraseToAnyPublisher()
    }

    func getAll(where predicate: NSPredicate?) -> [RepositoryObject] {
        var objects = self.realm.objects(RealmObject.self)

        if let predicate = predicate {
            objects = objects.filter(predicate)
        }
        return objects.compactMap { ($0).model as? RepositoryObject }
    }

    func insert(item: RepositoryObject) throws {
        try self.insertItem(item)
        self.subject.send(.insert)
    }

    func update(item: RepositoryObject) throws {
        try self.deleteItem(item)
        try self.insertItem(item)
        self.subject.send(.update)
    }

    func delete(item: RepositoryObject) throws {
        try self.deleteItem(item)
        self.subject.send(.delete)
    }

    private func insertItem(_ item: RepositoryObject) throws {
        try self.realm.write {
            realm.add(item.toStorable())
        }
    }

    private func deleteItem(_ item: RepositoryObject) throws {
        try self.realm.write {
            let predicate = NSPredicate(format: "uuid == %@",
                                        item.toStorable().uuid)
            if let productToDelete = realm.objects(RealmObject.self)
                .filter(predicate).first {
                realm.delete(productToDelete)
            }
        }
    }
}
