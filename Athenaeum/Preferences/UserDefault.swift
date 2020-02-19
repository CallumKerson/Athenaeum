//
//  UserDefault.swift
//  Athenaeum
//
//  Created by Callum Kerson on 19/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation

@propertyWrapper
struct UserDefault<T> {
  let key: Key
  let defaultValue: T
  
  var wrappedValue: T {
    get {
        return UserDefaults.standard.object(forKey: key.rawValue) as? T ?? defaultValue
    }
    set {
        UserDefaults.standard.set(newValue, forKey: key.rawValue)
    }
  }
}

struct Key: RawRepresentable {
    let rawValue: String
}
