#version 330

// Input vertex attributes (from vertex buffer objects)
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec3 vertexNormal;
in vec4 vertexColor;

// Input uniform values
uniform mat4 matModel;
uniform mat4 matView;
uniform mat4 matProjection;
uniform mat4 matNormal;

// Output vertex attributes to fragment shader
out vec2 fragTexCoord;
out vec3 fragNormal;

void main()
{
    // Calculate final vertex position
    vec4 worldPosition = matModel * vec4(vertexPosition, 1.0);
    gl_Position = matProjection * matView * worldPosition;
    
    // Calculate final normal in world space
    fragNormal = normalize(vec3(matNormal * vec4(vertexNormal, 0.0)));
    
    // Pass texCoord to fragment shader
    fragTexCoord = vertexTexCoord;
}

// Fragment shader
#ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
#else
precision mediump float;
#endif

in vec2 fragTexCoord;
in vec3 fragNormal;

uniform sampler2D texture0;
uniform sampler2D nightTexture;
uniform vec4 colDiffuse;
uniform vec3 lightDir;

out vec4 finalColor;

void main()
{
    float sunAmount = dot(normalize(fragNormal), normalize(lightDir));
    float dayMix = smoothstep(-0.08, 0.18, sunAmount);

    vec4 dayColor = texture(texture0, fragTexCoord)*colDiffuse;
    vec4 nightColor = texture(nightTexture, fragTexCoord);

    float daylight = 0.18 + max(sunAmount, 0.0)*0.92;
    vec3 litDay = dayColor.rgb*daylight;
    vec3 nightGlow = nightColor.rgb*1.35;

    finalColor = vec4(mix(nightGlow, litDay, dayMix), dayColor.a);
}
