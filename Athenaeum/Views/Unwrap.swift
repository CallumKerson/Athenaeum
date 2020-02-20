//
//  Unwrap.swift
//  Athenaeum
//
//  Created by Callum Kerson on 20/02/2020.
//  Copyright © 2020 Callum Kerson. All rights reserved.
//

import Foundation
import SwiftUI

struct Unwrap<Value, Content: View>: View {
    private let value: Value?
    private let contentProvider: (Value) -> Content

    init(_ value: Value?,
         @ViewBuilder content: @escaping (Value) -> Content) {
        self.value = value
        self.contentProvider = content
    }

    var body: some View {
        value.map(contentProvider)
    }
}
