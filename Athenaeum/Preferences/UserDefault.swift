/**
 UserDefault.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

@propertyWrapper
struct UserDefault<T> {
    let key: Key
    let defaultValue: T

    var wrappedValue: T {
        get {
            UserDefaults.standard
                .object(forKey: self.key.rawValue) as? T ?? self.defaultValue
        }
        set {
            UserDefaults.standard.set(newValue, forKey: self.key.rawValue)
        }
    }
}

struct Key: RawRepresentable {
    let rawValue: String
}
