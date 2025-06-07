### 
* create-modal 
	* create-dialog 
        * create-form
        
* edit-modal 
	 * edit-dialog 
        * edit-form

### Create Form

create button will hx-get of `/res/new` and target #create-modal, (swapping innerHTML by htmx's default)
close button will close the dialog.

### Create
create endpoint (POST /res) returns *table row by default, and oob swap of #create-modal to close the dialog
on form error, returns the form with error messages and set HX-Retarget, HX-Reswap headers to the form

### Edit form
edit button will hx-get of `/res/$id/edit` and target #edit-modal, (swapping innerHTML by htmx's default)
close button will close the dialog.

### Update
same as create. will send the <tr> on success, oob swap of #edit-modal, on error, returns the form with error messages
and set HX-Retarget, HX-Reswap headers

### Delete
will target the table row, and hx-swap of `delete`. will respond with 200 on success (204 does not swap)
    
### Internal Server Errors
HTMX will shows an alert box, no swapping will be done as htmx only swaps for 200/300 status codes


#### Index Page

```html
<button hx-get="/labels/new" hx-target="#create-modal">New</button>

<table>
    <thead>
        <th>Id</th>
        <th>Name</th>
        <th>Actions</th>
    </thead>
    <tbody id="contents">
    <tr id="row-1">
        <td>1</td>
        <td>Row1</td>
        <td>
            <button hx-get="/labels/1/edit" hx-target="#edit-modal">Edit</button>
            <button hx-delete="/labels/1" hx-target="#row-1" hx-swap="delete">Delete</button>
        </td>
    </tr>
    </tbody>
</table>

<div id="create-modal"></div>
<div id="edit-modal"></div>
<div id="details-modal"></div>

```

<details>
<summary>Create Form</summary>

```html
<dialog open id="create-dialog">
    <form hx-post="/labels"
        id="create-form"
        hx-target="#contents"
        hx-swap="afterbegin">
        <label>Name</label>
        <input type="text" name="name">

        <button type="submit">Create</button>
        <button type="button" onclick="$('#create-dialog').close()">Close</button>
    </form>
</dialog>
```

</details>


<details>
<summary>Create success</summary>

```html
    <tr id="row-9">
        <td>9</td>
        <td>helo world</td>
        <td>
            <button hx-get="/labels/9/edit" hx-target="#edit-modal">Edit</button>
            <button hx-delete="/labels/9" hx-target="#row-9" hx-swap="delete">Delete</button>
            <button>Details</button>
        </td>
    </tr>
    <div id="create-modal" hx-swap-oob="true"></div>
```

</details>


<details>
<summary>Create Errors</summary>
`HX-Retarget: #create-form`
`HX-Reswap: outerHTML`

```html
<form hx-post="/labels"
id="create-form"
hx-target="#contents"
hx-swap="afterbegin"
>
    
    <label>Name</label>
    <input type="text" name="name" aria-invalid="true" value="a">
    
        <small>string length must be between 4 and 20</small>
    

    <button type="submit">Create</button>
    <button type="button" onclick="$('#create-dialog').close()">Close</button>
</form>
```

</details>
 

<details>
<summary>Edit Form</summary>

```html
    <dialog open id="edit-dialog">
    <form hx-put="/labels/9"
    id="edit-form"
    hx-target="#row-9"
    hx-swap="outerHTML"
    >

    <label>Name</label>
    <input type="text" name="name" value="helo world">
    <button type="submit">Update</button>
    <button type="button" onclick="$('#edit-dialog').close()">Close</button>

    </form>
    </dialog>
```
</details>


<details>
<summary>Edit Success</summary>

```html
    <tr id="row-9">
        <td>9</td>
        <td>hello world</td>
        <td>
            <button hx-get="/labels/9/edit" hx-target="#edit-modal">Edit</button>
            <button hx-delete="/labels/9" hx-target="#row-9" hx-swap="delete">Delete</button>
            <button>Details</button>
        </td>
    </tr>
    <div id="edit-modal" hx-swap-oob="true"></div>

```
</details>

<details>
<summary>Edit Error</summary>
`HX-Retarget: #edit-form`
`HX-Reswap: innerHTML`

```html
    <form hx-put="/labels/9"
    id="edit-form"
    hx-target="#row-9"
    hx-swap="outerHTML"
    >

    <label>Name</label>
    <input type="text" name="name" value="a" aria-invalid="true">
    
        <small>string length must be between 4 and 20</small>
    
    <button type="submit">Update</button>
    <button type="button" onclick="$('#edit-dialog').close()">Close</button>

    </form>
```

</details>


<details>
<summary>Delete Success</summary>
200 Ok
</details>
