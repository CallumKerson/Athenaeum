/**
 EditGenreView.swift
 Copyright (c) 2023 Callum Kerr-Edwards
 */

import SwiftUI

struct EditGenreView: View {
    @ObservedObject var viewModel: EditGenreViewModel

    init(viewModel: inout EditGenreViewModel) {
        self.viewModel = viewModel
    }

    var body: some View {
        EnumPicker("Genre", selection: $viewModel.genreState)
    }
}

struct EditGenreView_Previews: PreviewProvider {
    static var emptyGenreViewModel: EditGenreViewModel = .init(initialGenre: nil)
    static var fantasyGenreViewModel: EditGenreViewModel = .init(initialGenre: .fantasy)

    static var previews: some View {
        Group {
            EditGenreView(viewModel: &emptyGenreViewModel)
            EditGenreView(viewModel: &fantasyGenreViewModel)
        }
    }
}

typealias Pickable = CaseIterable
    & Identifiable
    & Hashable
    & CustomStringConvertible

struct EnumPicker<Enum: Pickable, Label: View>: View {
    private let label: Label

    @Binding private var selection: Enum

    var body: some View {
        Picker(selection: $selection, label: label) {
            ForEach(Array(Enum.allCases)) { value in
                Text(value.description).tag(value)
            }
        }
    }

    init(selection: Binding<Enum>, label: Label) {
        self.label = label
        _selection = selection
    }
}

extension EnumPicker where Label == Text {
    init(_ titleKey: LocalizedStringKey, selection: Binding<Enum>) {
        self.label = Text(titleKey)
        _selection = selection
    }

    init<S: StringProtocol>(_ title: S, selection: Binding<Enum>) {
        self.label = Text(title)
        _selection = selection
    }
}
