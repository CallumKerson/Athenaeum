/**
 EditTextFieldViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

class EditTextFieldViewModel: ObservableObject {
    @Published var textState: String
    let initialTextValue: String
    let label: String

    init(label: String, initialText: String?) {
        self.label = label
        self.textState = initialText ?? ""
        self.initialTextValue = initialText ?? ""
    }

    func text() -> String? {
        if self.textState.isEmpty {
            return nil
        } else {
            return self.textState
        }
    }
}
