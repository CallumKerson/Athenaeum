/**
 Double+StringRepresentation.swift
 Copyright (c) 2020 Callum Kerr-Edwards
 */

import Foundation

extension Double {
    var asString: String {
        if self.truncatingRemainder(dividingBy: 1.0) == 0.0 {
            return "\(Int(self))"
        } else {
            return "\(self)"
        }
    }
}
