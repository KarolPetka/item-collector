import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {MainViewComponent} from './pages/main-view/main-view.component';
import {EditCollectionComponent} from "./pages/edit-collection/edit-collection.component";
import {NewCollectionComponent} from './pages/new-collection/new-collection.component';
import {NewItemComponent} from './pages/new-item/new-item.component';
import {EditItemComponent} from './pages/edit-item/edit-item.component';
import {LoginComponent} from "./pages/login/login.component";
import {SignupComponent} from "./pages/signup/signup.component";

const routes: Routes = [
  {path: "", redirectTo: "login", pathMatch: "full"},
  {path: "signup", component: SignupComponent},
  {path: "login", component: LoginComponent},
  {path: "collections", component: MainViewComponent},
  {path: "new-collection", component: NewCollectionComponent},
  {path: 'edit-collection/:collectionId', component: EditCollectionComponent},
  {path: "collections/:collectionId", component: MainViewComponent},
  {path: "collections/:collectionId/new-item", component: NewItemComponent},
  {path: 'collections/:collectionId/edit-item/:itemId', component: EditItemComponent},
  {path: '**', redirectTo: 'login', pathMatch: "full"},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class Routing {
}
