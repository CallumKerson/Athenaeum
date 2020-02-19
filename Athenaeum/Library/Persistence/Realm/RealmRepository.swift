/**
 GenericRepository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import RealmSwift
import Combine

class RealmRepository<RepositoryObject>: Repository
    where RepositoryObject: Entity,
RepositoryObject.StoreType: Object {
    typealias RealmObject = RepositoryObject.StoreType
    
    private let subject = PassthroughSubject<DatabaseAction, Never>()
    let publisher:AnyPublisher<DatabaseAction, Never>
    
    private let realm: Realm
    
    init() {
        realm = try! Realm()
        publisher = subject.eraseToAnyPublisher()
    }
    
    func getAll(where predicate: NSPredicate?) -> [RepositoryObject] {
        var objects = realm.objects(RealmObject.self)
        
        if let predicate = predicate {
            objects = objects.filter(predicate)
        }
        return objects.compactMap { ($0).model as? RepositoryObject }
    }
    
    func insert(item: RepositoryObject) throws {
        try insertItem(item)
        subject.send(.insert)
    }
    
    func update(item: RepositoryObject) throws {
        try deleteItem(item)
        try insertItem(item)
        subject.send(.update)
    }
    
    func delete(item: RepositoryObject) throws {
        try deleteItem(item)
        subject.send(.delete)
    }
    
    private func insertItem(_ item: RepositoryObject) throws {
        try realm.write {
            realm.add(item.toStorable())
        }
    }
    
    private func deleteItem(_ item: RepositoryObject) throws {
        try realm.write {
            let predicate = NSPredicate(format: "uuid == %@", item.toStorable().uuid)
            if let productToDelete = realm.objects(RealmObject.self)
                .filter(predicate).first {
                realm.delete(productToDelete)
            }
        }
    }
}
