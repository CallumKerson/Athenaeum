/**
 GenericRepository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Foundation
import RealmSwift

class RealmRepository<RepositoryObject>: Repository
    where RepositoryObject: Entity,
    RepositoryObject.StoreType: Object {
    typealias RealmObject = RepositoryObject.StoreType

    private let realm: Realm

    init() {
        realm = try! Realm()
    }

    func getAll(where predicate: NSPredicate?) -> [RepositoryObject] {
        var objects = realm.objects(RealmObject.self)

        if let predicate = predicate {
            objects = objects.filter(predicate)
        }
        return objects.compactMap { ($0).model as? RepositoryObject }
    }

    func insert(item: RepositoryObject) throws {
        try realm.write {
            realm.add(item.toStorable())
        }
    }

    func update(item: RepositoryObject) throws {
        try delete(item: item)
        try insert(item: item)
    }

    func delete(item: RepositoryObject) throws {
        try realm.write {
            let predicate = NSPredicate(format: "uuid == %@", item.toStorable().uuid)
            if let productToDelete = realm.objects(RealmObject.self)
                .filter(predicate).first {
                realm.delete(productToDelete)
            }
        }
    }
}
