import {Injectable} from '@angular/core';
import {Repository} from './repository';

@Injectable({
  providedIn: 'root'
})
export class Service {

  constructor(private repository: Repository) {
  }

  signup(username: string, password: string) {
    const payload = JSON.stringify({ username, password });
    return this.repository.postWithoutHeaders("signup", payload);
  }

  login(username: string, password: string) {
    const payload = JSON.stringify({ username, password });
    return this.repository.postWithoutHeaders("login", payload);
  }

  getCollections(token: string) {
    return this.repository.get('collections', token);
  }

  createCollection(title: string, token: string) {
    return this.repository.post("collections", {title}, token);
  }

  updateCollection(id: string, title: string, token: string) {
    return this.repository.put(`collections/${id}`, {title}, token);
  }

  deleteCollection(id: string, token: string) {
    return this.repository.delete(`collections/${id}`, token);
  }

  getItems(collectionId: string, token: string) {
    return this.repository.get(`collections/${collectionId}/items`, token);
  }

  createItem(item: any, collectionId: string, token: string) {
    return this.repository.postJson(`collections/${collectionId}/items`, item, token);
  }

  updateItem(item: any, collectionId: string, itemId: string, token: string) {
    return this.repository.put(`collections/${collectionId}/items/${itemId}`, item, token);
  }

  deleteItem(collectionId: string, itemId: string, token: string) {
    return this.repository.delete(`collections/${collectionId}/items/${itemId}`, token);
  }
}
