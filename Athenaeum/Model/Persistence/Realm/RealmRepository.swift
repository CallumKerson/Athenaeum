/**
 RealmRepository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation
import RealmSwift

class RealmRepository<RepositoryObject>: Repository,
    ObservableObject where RepositoryObject: Entity,
    RepositoryObject.StoreType: Object {
    typealias RealmObject = RepositoryObject.StoreType

    @Published var items = [RepositoryObject]()

    private var token: NotificationToken?

    init() {
        log.debug("Initalising Realm Repository")
        self.reloadItems()
        self.token = try! Realm().objects(RealmObject.self).observe { _ in
            self.reloadItems()
        }
    }

    deinit {
        token?.invalidate()
    }

    private func reloadItems() {
        self.items = try! Realm()
            .objects(RealmObject.self)
            .compactMap { ($0).model as? RepositoryObject }
    }

    func insert(item: RepositoryObject) throws {
        let realm = try! Realm()
        try realm.write {
            realm.add(item.toStorable())
        }
    }

    func update(item: RepositoryObject) throws {
        try self.delete(item: item)
        try self.insert(item: item)
    }

    func delete(item: RepositoryObject) throws {
        let realm = try! Realm()
        try realm.write {
            let predicate = NSPredicate(format: "uuid == %@",
                                        item.toStorable().uuid)
            if let productToDelete = realm.objects(RealmObject.self)
                .filter(predicate).first {
                realm.delete(productToDelete)
            }
        }
    }
}
