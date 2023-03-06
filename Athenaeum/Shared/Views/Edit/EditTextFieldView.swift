/**
 EditTextFieldView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct EditTextFieldView: View {
    @ObservedObject var viewModel: EditTextFieldViewModel

    init(viewModel: inout EditTextFieldViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        HStack {
            Text(viewModel.label)
            Spacer()
            TextField(viewModel.initialTextValue, text: $viewModel.textState)
        }
    }
}

struct EditTextFieldView_Previews: PreviewProvider {
    static var titleViewModel: EditTextFieldViewModel = .init(
        label: "Title",
        initialText: "A Book Title"
    )

    static var previews: some View {
        EditTextFieldView(viewModel: &titleViewModel)
    }
}
