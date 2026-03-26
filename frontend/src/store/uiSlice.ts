import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { Article } from "../types";

export type View =
  | "dashboard"
  | "articles"
  | "article"
  | "edit"
  | "settings"
  | "reader-preview";

interface UiState {
  view: View;
  selectedArticleId: string | null;
  focusSectionId: string | null;
  previewArticle: Article | null;
  activeModal: string | null;
  activeMenu: string | null;
}

const initialState: UiState = {
  view: "dashboard",
  selectedArticleId: null,
  focusSectionId: null,
  previewArticle: null,
  activeModal: null,
  activeMenu: null,
};

export const uiSlice = createSlice({
  name: "ui",
  initialState,
  reducers: {
    setView(state, action: PayloadAction<View>) {
      state.view = action.payload;
    },
    selectArticle(state, action: PayloadAction<string>) {
      state.selectedArticleId = action.payload;
    },
    setFocusSection(state, action: PayloadAction<string | null>) {
      state.focusSectionId = action.payload;
    },
    setPreviewArticle(state, action: PayloadAction<Article | null>) {
      state.previewArticle = action.payload;
    },
    goBack(state) {
      state.view = "dashboard";
      state.selectedArticleId = null;
      state.focusSectionId = null;
      state.previewArticle = null;
    },
    openModal(state, action: PayloadAction<string>) {
      state.activeModal = action.payload;
    },
    closeModal(state) {
      state.activeModal = null;
    },
    setActiveMenu(state, action: PayloadAction<string | null>) {
      state.activeMenu = action.payload;
    },
  },
});

export const {
  setView,
  selectArticle,
  setFocusSection,
  setPreviewArticle,
  goBack,
  openModal,
  closeModal,
  setActiveMenu,
} = uiSlice.actions;
