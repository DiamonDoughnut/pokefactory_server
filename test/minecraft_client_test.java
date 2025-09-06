// Simple Java HTTP client to test API calls (mimics Minecraft mod behavior)
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.net.URI;
import java.time.Duration;

public class MinecraftClientTest {
    private static final String API_URL = "http://localhost:8080";
    private static String serverToken = "";
    
    public static void main(String[] args) {
        System.out.println("=== Minecraft Client API Test ===");
        
        try {
            // 1. Authenticate as server
            authenticateServer();
            
            // 2. Create player (simulates player joining server)
            createPlayer("550e8400-e29b-41d4-a716-446655440000", "TestPlayer");
            
            // 3. Simulate catching Pokémon
            catchPokemon("550e8400-e29b-41d4-a716-446655440000", 25, "catch"); // Pikachu
            catchPokemon("550e8400-e29b-41d4-a716-446655440000", 1, "catch");  // Bulbasaur
            catchPokemon("550e8400-e29b-41d4-a716-446655440000", 150, "see");  // Mewtwo (seen only)
            
            // 4. Get player summary
            getPlayerSummary("550e8400-e29b-41d4-a716-446655440000");
            
            System.out.println("\\n=== Test Complete ===");
            
        } catch (Exception e) {
            System.err.println("Test failed: " + e.getMessage());
        }
    }
    
    private static void authenticateServer() throws Exception {
        String json = "{\"server_id\":\"test-server-1\",\"server_key\":\"your-jwt-secret\"}";
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(API_URL + "/api/v1/server/auth"))
            .header("Content-Type", "application/json")
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .timeout(Duration.ofSeconds(10))
            .build();
            
        HttpResponse<String> response = HttpClient.newHttpClient().send(request, 
            HttpResponse.BodyHandlers.ofString());
            
        System.out.println("Server Auth: " + response.statusCode());
        
        // Extract token (simplified - would use JSON parser in real implementation)
        String body = response.body();
        if (body.contains("\"token\":")) {
            serverToken = body.split("\"token\":\"")[1].split("\"")[0];
            System.out.println("Got server token: " + serverToken.substring(0, 20) + "...");
        }
    }
    
    private static void createPlayer(String uuid, String username) throws Exception {
        String json = String.format("{\"player_uuid\":\"%s\",\"username\":\"%s\"}", uuid, username);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(API_URL + "/api/v1/server/player/create"))
            .header("Content-Type", "application/json")
            .header("Authorization", "Bearer " + serverToken)
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
            
        HttpResponse<String> response = HttpClient.newHttpClient().send(request, 
            HttpResponse.BodyHandlers.ofString());
            
        System.out.println("Create Player: " + response.statusCode() + " - " + response.body());
    }
    
    private static void catchPokemon(String uuid, int nationalId, String action) throws Exception {
        String json = String.format("{\"player_uuid\":\"%s\",\"national_id\":%d,\"action\":\"%s\"}", 
            uuid, nationalId, action);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(API_URL + "/api/v1/server/pokedex/update"))
            .header("Content-Type", "application/json")
            .header("Authorization", "Bearer " + serverToken)
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
            
        HttpResponse<String> response = HttpClient.newHttpClient().send(request, 
            HttpResponse.BodyHandlers.ofString());
            
        System.out.println("Pokémon " + nationalId + " " + action + ": " + response.statusCode());
    }
    
    private static void getPlayerSummary(String uuid) throws Exception {
        String json = String.format("{\"player_uuid\":\"%s\"}", uuid);
        
        HttpRequest request = HttpRequest.newBuilder()
            .uri(URI.create(API_URL + "/api/v1/server/pokedex/summary"))
            .header("Content-Type", "application/json")
            .header("Authorization", "Bearer " + serverToken)
            .POST(HttpRequest.BodyPublishers.ofString(json))
            .build();
            
        HttpResponse<String> response = HttpClient.newHttpClient().send(request, 
            HttpResponse.BodyHandlers.ofString());
            
        System.out.println("Player Summary: " + response.body());
    }
}