function updateCharCount() {
    const textarea = document.getElementById('messageContent');
    const charCount = document.getElementById('charCount');
    const maxLength = 200;
    const currentLength = textarea.value.length;

    charCount.textContent = `${currentLength} / ${maxLength}`;
}

document.getElementById('infoToggle').addEventListener('click', function() {
    const infoList = document.getElementById('infoList');
    const triangle = document.getElementById('triangle');
    
    if (infoList.style.maxHeight === '0px' || infoList.style.maxHeight === '') {
        infoList.style.maxHeight = '200px';
        triangle.innerHTML = '▲';
    } else {
        infoList.style.maxHeight = '0px';
        triangle.innerHTML = '▼';
    }
});
