//
//  Persistable.swift
//  Athenaeum
//
//  Created by Callum Kerson on 17/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation

protocol Repository {
    associatedtype EntityObject: Entity
    
    func getAll(where predicate: NSPredicate?) -> [EntityObject]
    func insert(item: EntityObject) throws
    func update(item: EntityObject) throws
    func delete(item: EntityObject) throws
}

extension Repository {
    func getAll() -> [EntityObject] {
        return getAll(where: nil)
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
