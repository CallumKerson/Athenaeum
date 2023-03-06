/**
 EditMultilineTextFieldViewModel.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Foundation

class EditMultilineTextFieldViewModel: ObservableObject {
    @Published var multilineTextState: String
    let sectionLabel: String

    init(sectionLabel: String, initialMultilineText: String?) {
        self.sectionLabel = sectionLabel
        self.multilineTextState = initialMultilineText ?? ""
    }

    func multilineText() -> String? {
        if self.multilineTextState.isEmpty {
            return nil
        } else {
            return self.multilineTextState
        }
    }
}
