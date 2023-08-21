<script>

async function handleSubmit(e) {
    e.preventDefault();

    const url = new URL(e.target.action)
    const formData = new FormData(e.target);

    const res = await (await fetch(url, {
        method: "POST",
        body: formData
    })).text()


    console.log({res})
}

</script>

<div class="h-80 p-6">
    <form on:submit={handleSubmit} action="/file" method="POST" enctype="multipart/form-data">
        <label for="file">Files</label>
        <input name="files" multiple="true" type="file" />
        <input class="rounded bg-blue-500 py-2 px-4 text-white" type="submit" />
    </form>
</div>
