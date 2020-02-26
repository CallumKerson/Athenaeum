/**
 Repository.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Combine
import Foundation

protocol Repository: ObservableObject {
    associatedtype EntityObject: Entity

    var items: [EntityObject] { get set }

    func insert(item: EntityObject) throws
    func update(item: EntityObject) throws
    func delete(item: EntityObject) throws
    func objectExists(item: EntityObject) throws -> Bool
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
