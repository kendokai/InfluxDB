import { useState, useEffect } from 'react';

async function fetchUID(setUID) {
  console.log("Fetching UID froms/UID...");
  try {
    const response = await fetch(`/get-uid`);
    if (!response.ok) {
        throw new Error('Failed to fetch uid')
    }
    const data = await response.json();
    console.log(data)
    setUID(data.uid);
} catch (error) {
    console.error('Error fetching uid:', error);
}
}

export default function GrafanaPanel({timeRange}) {
  const [UID, setUID] = useState('')
  const grafanaURL = "http://localhost:3000/d-solo/"
  const grafanaQueryText = "/trial?orgId=1&panelId=1"
  useEffect(() => {
    fetchUID(setUID);
  }, []);
  
  return (
    <div>
      <iframe
        src={grafanaURL + UID + grafanaQueryText + "&from=" + timeRange.from + "&to=" + timeRange.to}
        width="100%"
        height="500"
        frameBorder="0"
        title="Grafana Time Series"
        key={timeRange.state}
      ></iframe>
    </div>
  );
};