/**
 EditMultilineTextFieldView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import Combine
import SwiftUI

struct EditMultilineTextFieldView: View {
    @ObservedObject var viewModel: EditMultilineTextFieldViewModel

    init(viewModel: inout EditMultilineTextFieldViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        Section(header: Text("Summary").font(.headline)) {
            TextEditor(text: $viewModel.multilineTextState)
                .shadow(radius: 1)
                .frame(minWidth: 200, minHeight: 100)
        }
    }
}

struct SummarySection_Previews: PreviewProvider {
    static var summaryViewModel: EditMultilineTextFieldViewModel =
        .init(
            sectionLabel: "Summary",
            initialMultilineText: "Prexisting summary of the book.\n\n It can be multiple lines."
        )

    static var previews: some View {
        EditMultilineTextFieldView(viewModel: &summaryViewModel)
    }
}
