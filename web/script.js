async function getOrder() {
    const orderId = document.getElementById('orderId').value;
    if (!orderId) return;

    try {
        const response = await fetch(`http://localhost:8081/order/${orderId}`);
        if (!response.ok) throw new Error('Order not found');
        
        const data = await response.json();
        document.getElementById('result').textContent = JSON.stringify(data, null, 2);
    } catch (error) {
        document.getElementById('result').textContent = `Error: ${error.message}`;
    }
}