/**
 Repository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 Licensed under the MIT license.
 */

import Combine
import Foundation

protocol Repository {
    associatedtype EntityObject: Entity

    var publisher: AnyPublisher<DatabaseAction, Never> { get }

    func getAll(where predicate: NSPredicate?) -> [EntityObject]
    func insert(item: EntityObject) throws
    func update(item: EntityObject) throws
    func delete(item: EntityObject) throws
}

extension Repository {
    func getAll() -> [EntityObject] {
        getAll(where: nil)
    }
}

public protocol Entity {
    associatedtype StoreType: Storable

    func toStorable() -> StoreType
}

public protocol Storable {
    associatedtype EntityObject: Entity

    var model: EntityObject { get }
    var uuid: String { get }
}

public enum DatabaseAction: String {
    case insert
    case update
    case delete
}
